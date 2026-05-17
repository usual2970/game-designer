package threeos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validPolicyToken() *FilePolicyTokenResp {
	return &FilePolicyTokenResp{
		Policy:           "test-policy-base64",
		SecurityToken:    "test-security-token",
		SignatureVersion: "4.0",
		Credential:       "test-cred",
		Date:             "20250626T000000Z",
		Signature:        "test-signature",
		Host:             "https://oss.example.com",
		Dir:              "uploads/2025/06/26/",
	}
}

func TestUploadFile_PackageUpload(t *testing.T) {
	var receivedKey, receivedPolicy, receivedSig string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		r.ParseMultipartForm(10 << 20)
		receivedKey = r.FormValue("key")
		receivedPolicy = r.FormValue("policy")
		receivedSig = r.FormValue("signature")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("fake zip content"), 0644)

	policy := validPolicyToken()
	policy.Host = server.URL

	uploader := NewOSSUploader(server.Client())
	result, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(result.ObjectURL, ".zip") {
		t.Errorf("expected URL ending in .zip, got %s", result.ObjectURL)
	}
	if !strings.HasPrefix(receivedKey, policy.Dir) {
		t.Errorf("expected key starting with %s, got %s", policy.Dir, receivedKey)
	}
	if receivedPolicy != "test-policy-base64" {
		t.Errorf("expected policy in form, got %s", receivedPolicy)
	}
	if receivedSig != "test-signature" {
		t.Errorf("expected signature in form, got %s", receivedSig)
	}
}

func TestUploadFile_SQLUpload(t *testing.T) {
	var receivedKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		receivedKey = r.FormValue("key")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "init.sql")
	os.WriteFile(tmpFile, []byte("CREATE TABLE test;"), 0644)

	policy := validPolicyToken()
	policy.Host = server.URL

	uploader := NewOSSUploader(server.Client())
	result, err := uploader.UploadFile(context.Background(), tmpFile, policy, "sql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(result.ObjectURL, ".sql") {
		t.Errorf("expected URL ending in .sql, got %s", result.ObjectURL)
	}
	if !strings.HasSuffix(receivedKey, ".sql") {
		t.Errorf("expected key ending in .sql, got %s", receivedKey)
	}
}

func TestUploadFile_EmptyPathIgnored(t *testing.T) {
	uploader := NewOSSUploader(http.DefaultClient)
	result, err := uploader.UploadFile(context.Background(), "", validPolicyToken(), "optional")
	if err != nil {
		t.Fatalf("expected nil error for empty path, got %v", err)
	}
	if result != nil {
		t.Error("expected nil result for empty path")
	}
}

func TestUploadFile_FileNotFound(t *testing.T) {
	uploader := NewOSSUploader(http.DefaultClient)
	_, err := uploader.UploadFile(context.Background(), "/nonexistent/file.zip", validPolicyToken(), "package")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestUploadFile_MissingPolicyHost(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	policy := validPolicyToken()
	policy.Host = ""

	uploader := NewOSSUploader(http.DefaultClient)
	_, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err == nil {
		t.Fatal("expected error for missing host")
	}
	if !strings.Contains(err.Error(), "missing host") {
		t.Errorf("expected 'missing host' error, got: %v", err)
	}
}

func TestUploadFile_MissingPolicyDir(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	policy := validPolicyToken()
	policy.Dir = ""

	uploader := NewOSSUploader(http.DefaultClient)
	_, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err == nil {
		t.Fatal("expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "missing host or dir") {
		t.Errorf("expected 'missing host or dir' error, got: %v", err)
	}
}

func TestUploadFile_MissingPolicySignature(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	policy := validPolicyToken()
	policy.Policy = ""
	policy.Signature = ""

	uploader := NewOSSUploader(http.DefaultClient)
	_, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err == nil {
		t.Fatal("expected error for missing policy/signature")
	}
	if !strings.Contains(err.Error(), "missing policy or signature") {
		t.Errorf("expected 'missing policy or signature' error, got: %v", err)
	}
}

func TestUploadFile_OSSRejectsUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("AccessDenied"))
	}))
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	policy := validPolicyToken()
	policy.Host = server.URL

	uploader := NewOSSUploader(server.Client())
	_, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err == nil {
		t.Fatal("expected error for OSS rejection")
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Errorf("expected upload failure error, got: %v", err)
	}
}

func TestUploadFile_MultipartFieldsPresent(t *testing.T) {
	var formFields map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		formFields = map[string]string{
			"policy":   r.FormValue("policy"),
			"security": r.FormValue("x-oss-security-token"),
			"cred":     r.FormValue("x-oss-credential"),
			"date":     r.FormValue("x-oss-date"),
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	policy := validPolicyToken()
	policy.Host = server.URL

	uploader := NewOSSUploader(server.Client())
	_, err := uploader.UploadFile(context.Background(), tmpFile, policy, "package")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if formFields["policy"] != "test-policy-base64" {
		t.Errorf("expected policy field, got %s", formFields["policy"])
	}
	if formFields["security"] != "test-security-token" {
		t.Errorf("expected security token field, got %s", formFields["security"])
	}
	if formFields["cred"] != "test-cred" {
		t.Errorf("expected credential field, got %s", formFields["cred"])
	}
	if formFields["date"] != "20250626T000000Z" {
		t.Errorf("expected date field, got %s", formFields["date"])
	}
}

func TestUniqueFilename(t *testing.T) {
	name := uniqueFilename("/path/to/my-game.zip")
	if !strings.HasSuffix(name, ".zip") {
		t.Errorf("expected .zip extension, got %s", name)
	}
	if len(name) < 10 {
		t.Errorf("expected longer unique name, got %s", name)
	}
}

func TestUniqueFilename_NoExtension(t *testing.T) {
	name := uniqueFilename("/path/to/README")
	if strings.Contains(name, "..") {
		t.Errorf("expected no double dots, got %s", name)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with spaces", "withspa"},
		{"special!@#chars", "special"},
		{"123numbers", "123numbe"},
	}
	for _, tt := range tests {
		result := sanitizeFilename(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
