package threeos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) Login(ctx context.Context, identifier, password string) (*AuthLoginResp, error) {
	reqBody := AuthLoginReq{
		Identifier: identifier,
		Type:       "password",
		Data:       password,
	}
	var resp AuthLoginResp
	if err := c.doPost(ctx, "/common/v1/auth/login", reqBody, &resp, false); err != nil {
		return nil, err
	}
	c.token = resp.AccessToken
	return &resp, nil
}

func (c *Client) GetPolicyToken(ctx context.Context) (*FilePolicyTokenResp, error) {
	var resp FilePolicyTokenResp
	if err := c.doGet(ctx, "/developer/v1/file/policy-token", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ListGames(ctx context.Context, page, pageSize int) (*GameListResp, error) {
	path := fmt.Sprintf("/developer/v1/game?page=%d&pageSize=%d", page, pageSize)
	var resp GameListResp
	if err := c.doGet(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateWithVersion(ctx context.Context, req *GameCreateWithVersionReq) (*GameInfoResp, error) {
	var resp GameInfoResp
	if err := c.doPost(ctx, "/developer/v1/game/create-with-version", req, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateGame(ctx context.Context, gameURI string, req *GameUpdateReq) (*GameInfoResp, error) {
	path := "/developer/v1/game/" + url.PathEscape(gameURI)
	var resp GameInfoResp
	if err := c.doPut(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateWithVersion(ctx context.Context, req *GameUpdateWithVersionReq) (*GameInfoResp, error) {
	var resp GameInfoResp
	if err := c.doPost(ctx, "/developer/v1/game/update-with-version", req, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ApplyReview(ctx context.Context, reviewURI string) error {
	reqBody := GameReviewApplyReq{URI: reviewURI}
	return c.doPost(ctx, "/developer/v1/game/apply-review", reqBody, nil, true)
}

func (c *Client) doGet(ctx context.Context, path string, result interface{}) error {
	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return &ClientError{Endpoint: path, Message: fmt.Sprintf("build request: %v", err)}
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.doRequest(req, path, result)
}

func (c *Client) doPost(ctx context.Context, path string, body interface{}, result interface{}, auth bool) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return &ClientError{Endpoint: path, Message: fmt.Sprintf("marshal body: %v", err)}
	}
	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return &ClientError{Endpoint: path, Message: fmt.Sprintf("build request: %v", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	if auth && c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.doRequest(req, path, result)
}

func (c *Client) doPut(ctx context.Context, path string, body interface{}, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return &ClientError{Endpoint: path, Message: fmt.Sprintf("marshal body: %v", err)}
	}
	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return &ClientError{Endpoint: path, Message: fmt.Sprintf("build request: %v", err)}
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.doRequest(req, path, result)
}

func (c *Client) doRequest(req *http.Request, endpoint string, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ClientError{Endpoint: endpoint, Message: fmt.Sprintf("request failed: %v", err)}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return &ClientError{Endpoint: endpoint, Message: fmt.Sprintf("read response: %v", err)}
	}

	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return &ClientError{Endpoint: endpoint, Message: fmt.Sprintf("decode response: %v", err)}
	}
	if !apiResp.IsSuccess() {
		return &ClientError{Endpoint: endpoint, Message: apiResp.Message, Code: apiResp.Code}
	}
	if result != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return &ClientError{Endpoint: endpoint, Message: fmt.Sprintf("decode data: %v", err)}
		}
	}
	return nil
}

func BuildObjectURL(host, dir, filename string) string {
	return host + "/" + dir + filename
}
