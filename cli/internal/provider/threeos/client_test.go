package threeos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func writeAPIResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := APIResponse{Code: code, Message: message}
	if data != nil {
		b, _ := json.Marshal(data)
		resp.Data = b
	}
	json.NewEncoder(w).Encode(resp)
}

func TestClient_Login(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/common/v1/auth/login" {
			t.Errorf("expected path /common/v1/auth/login, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req AuthLoginReq
		json.NewDecoder(r.Body).Decode(&req)
		if req.Identifier != "test@example.com" {
			t.Errorf("expected identifier=test@example.com, got %s", req.Identifier)
		}
		if req.Type != "password" {
			t.Errorf("expected type=password, got %s", req.Type)
		}

		writeAPIResponse(w, 0, "success", &AuthLoginResp{
			AccessToken: "test-token-123",
			ExpiresAt:   1234567890,
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	resp, err := client.Login(context.Background(), "test@example.com", "testpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "test-token-123" {
		t.Errorf("expected token test-token-123, got %s", resp.AccessToken)
	}
	if client.Token() != "test-token-123" {
		t.Errorf("client did not store token after login")
	}
}

func TestClient_Login_StoresBearerToken(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/common/v1/auth/login" {
			writeAPIResponse(w, 0, "success", &AuthLoginResp{AccessToken: "stored-token"})
			return
		}
		authHeader = r.Header.Get("Authorization")
		writeAPIResponse(w, 0, "success", &GameListResp{})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.Login(context.Background(), "user", "pass")
	client.ListGames(context.Background(), 1, 10)

	if authHeader != "Bearer stored-token" {
		t.Errorf("expected Bearer stored-token, got %s", authHeader)
	}
}

func TestClient_GetPolicyToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developer/v1/file/policy-token" {
			t.Errorf("expected path /developer/v1/file/policy-token, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		writeAPIResponse(w, 0, "success", &FilePolicyTokenResp{
			Policy:        "test-policy",
			SecurityToken: "test-security-token",
			Host:          "https://oss.example.com",
			Dir:           "uploads/2025/06/26/",
			Signature:     "test-sig",
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	resp, err := client.GetPolicyToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Host != "https://oss.example.com" {
		t.Errorf("expected host, got %s", resp.Host)
	}
	if resp.Dir != "uploads/2025/06/26/" {
		t.Errorf("expected dir, got %s", resp.Dir)
	}
}

func TestClient_ListGames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developer/v1/game" {
			t.Errorf("expected path /developer/v1/game, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("pageSize") != "5" {
			t.Errorf("expected page=2&pageSize=5, got %s", r.URL.RawQuery)
		}
		writeAPIResponse(w, 0, "success", &GameListResp{
			Page:       2,
			PageSize:   5,
			TotalPages: 3,
			TotalCount: 12,
			Data: []GameInfoResp{
				{URI: "game1", Name: "Game One"},
				{URI: "game2", Name: "Game Two"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	resp, err := client.ListGames(context.Background(), 2, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 12 {
		t.Errorf("expected totalCount=12, got %d", resp.TotalCount)
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 games, got %d", len(resp.Data))
	}
	if resp.Data[0].URI != "game1" {
		t.Errorf("expected first game uri=game1, got %s", resp.Data[0].URI)
	}
}

func TestClient_CreateWithVersion(t *testing.T) {
	var receivedBody GameCreateWithVersionReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developer/v1/game/create-with-version" {
			t.Errorf("expected path /developer/v1/game/create-with-version, got %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		writeAPIResponse(w, 0, "success", &GameInfoResp{
			URI:  "newgame123",
			Name: "Test Game",
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	req := &GameCreateWithVersionReq{
		Name:        "Test Game",
		Description: "A test game",
		Version: GameVersionCreateReq{
			Version:   "1.0.0",
			ChangeLog: "Initial release",
			FileUrl:   "https://oss.example.com/game.zip",
		},
	}
	resp, err := client.CreateWithVersion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.URI != "newgame123" {
		t.Errorf("expected uri=newgame123, got %s", resp.URI)
	}
	if receivedBody.Name != "Test Game" {
		t.Errorf("expected name=Test Game, got %s", receivedBody.Name)
	}
	if receivedBody.Version.FileUrl != "https://oss.example.com/game.zip" {
		t.Errorf("expected fileUrl in request")
	}
}

func TestClient_UpdateGame(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/developer/v1/game/some-game-uri"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeAPIResponse(w, 0, "success", &GameInfoResp{
			URI:  "some-game-uri",
			Name: "Updated Game",
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	req := &GameUpdateReq{
		URI:         "some-game-uri",
		Description: "Updated description",
	}
	resp, err := client.UpdateGame(context.Background(), "some-game-uri", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.URI != "some-game-uri" {
		t.Errorf("expected uri=some-game-uri, got %s", resp.URI)
	}
}

func TestClient_UpdateWithVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developer/v1/game/update-with-version" {
			t.Errorf("expected path /developer/v1/game/update-with-version, got %s", r.URL.Path)
		}
		writeAPIResponse(w, 0, "success", &GameInfoResp{
			URI:  "existinggame",
			Name: "Updated Version",
		})
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	req := &GameUpdateWithVersionReq{
		URI:   "existinggame",
		Name:  "Updated Version",
		Version: GameVersionCreateReq{
			Version:   "1.1.0",
			ChangeLog: "Bug fixes",
			FileUrl:   "https://oss.example.com/game-v2.zip",
		},
	}
	resp, err := client.UpdateWithVersion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.URI != "existinggame" {
		t.Errorf("expected uri=existinggame, got %s", resp.URI)
	}
}

func TestClient_ApplyReview(t *testing.T) {
	var receivedBody GameReviewApplyReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developer/v1/game/apply-review" {
			t.Errorf("expected path /developer/v1/game/apply-review, got %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		writeAPIResponse(w, 0, "success", nil)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	err := client.ApplyReview(context.Background(), "review-uri-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedBody.URI != "review-uri-123" {
		t.Errorf("expected uri=review-uri-123, got %s", receivedBody.URI)
	}
}

func TestClient_ApplyReview_NonZeroCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIResponse(w, 100, "review already submitted", nil)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	err := client.ApplyReview(context.Background(), "review-uri")
	if err == nil {
		t.Fatal("expected error for non-zero code")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("expected *ClientError, got %T", err)
	}
	if clientErr.Code != 100 {
		t.Errorf("expected code=100, got %d", clientErr.Code)
	}
	if !strings.Contains(clientErr.Message, "review already submitted") {
		t.Errorf("expected message to contain error text, got %s", clientErr.Message)
	}
}

func TestClient_NonZeroCodeApplicationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIResponse(w, 200, "unauthorized", nil)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	_, err := client.ListGames(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error for non-zero envelope code")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("expected *ClientError, got %T", err)
	}
	if clientErr.Code != 200 {
		t.Errorf("expected code=200, got %d", clientErr.Code)
	}
}

func TestClient_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	client.SetToken("test-token")
	_, err := client.ListGames(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	if !strings.Contains(err.Error(), "decode response") {
		t.Errorf("expected decode error, got: %v", err)
	}
}

func TestClient_NetworkError(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost:0")
	_, err := client.ListGames(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
}

func TestClient_LoginErrorResponseNoSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIResponse(w, 100, "invalid credentials", nil)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL)
	_, err := client.Login(context.Background(), "user@test.com", "secret-pass-123")
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "secret-pass-123") {
		t.Error("password leaked into error message")
	}
	if strings.Contains(err.Error(), "user@test.com") {
		t.Error("identifier leaked into error message")
	}
}

func TestBuildObjectURL(t *testing.T) {
	url := BuildObjectURL("https://oss.example.com", "uploads/2025/06/", "game.zip")
	expected := "https://oss.example.com/uploads/2025/06/game.zip"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}
