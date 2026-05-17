package threeos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OSSUploader struct {
	httpClient *http.Client
}

func NewOSSUploader(httpClient *http.Client) *OSSUploader {
	return &OSSUploader{httpClient: httpClient}
}

type UploadResult struct {
	Label     string
	LocalPath string
	ObjectURL string
	ObjectKey string
}

func (u *OSSUploader) UploadFile(ctx context.Context, localPath string, policy *FilePolicyTokenResp, label string) (*UploadResult, error) {
	if localPath == "" {
		return nil, nil
	}

	if policy.Host == "" || policy.Dir == "" {
		return nil, fmt.Errorf("upload: policy token missing host or dir")
	}
	if policy.Policy == "" || policy.Signature == "" {
		return nil, fmt.Errorf("upload: policy token missing policy or signature")
	}

	objectName := policy.Dir + uniqueFilename(localPath)
	objectURL := policy.Host + "/" + objectName

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := map[string]string{
		"policy":                  policy.Policy,
		"x-oss-security-token":    policy.SecurityToken,
		"x-oss-signature-version": policy.SignatureVersion,
		"x-oss-credential":        policy.Credential,
		"x-oss-date":              policy.Date,
		"signature":               policy.Signature,
		"key":                     objectName,
	}
	for key, val := range fields {
		if err := writer.WriteField(key, val); err != nil {
			return nil, fmt.Errorf("upload: write field %s: %w", key, err)
		}
	}

	file, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("upload: %s file not found: %w", label, err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(localPath))
	if err != nil {
		return nil, fmt.Errorf("upload: create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("upload: copy file content: %w", err)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, policy.Host, body)
	if err != nil {
		return nil, fmt.Errorf("upload: build request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload: %s request failed: %w", label, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload: %s failed (status %d): %s", label, resp.StatusCode, string(respBody))
	}

	return &UploadResult{
		Label:     label,
		LocalPath: localPath,
		ObjectURL: objectURL,
		ObjectKey: objectName,
	}, nil
}

func uniqueFilename(originalPath string) string {
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(filepath.Base(originalPath), ext)
	now := time.Now()
	timestamp := now.Format("0601021504")
	shortRand := fmt.Sprintf("%04d", now.UnixNano()%10000)
	return fmt.Sprintf("%s%s%s%s", timestamp, shortRand, sanitizeFilename(base), ext)
}

func sanitizeFilename(name string) string {
	runes := []rune(name)
	if len(runes) > 8 {
		runes = runes[:8]
	}
	result := strings.Builder{}
	for _, r := range runes {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}
	return result.String()
}
