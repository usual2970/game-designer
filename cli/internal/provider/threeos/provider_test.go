package threeos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/game-designer-cli/internal/provider"
)

type testEnv struct {
	apiServer  *httptest.Server
	ossServer  *httptest.Server
	prov       *ThreeOSProvider
	tmpDir     string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	ossUploads := 0
	ossServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ossUploads++
		w.WriteHeader(http.StatusNoContent)
	}))

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/common/v1/auth/login":
			writeAPIResponse(w, 0, "success", &AuthLoginResp{AccessToken: "test-token"})
		case r.URL.Path == "/developer/v1/file/policy-token":
			writeAPIResponse(w, 0, "success", &FilePolicyTokenResp{
				Policy:    "test-policy",
				Signature: "test-sig",
				Host:      ossServer.URL,
				Dir:       "uploads/test/",
			})
		case r.URL.Path == "/developer/v1/game" && r.Method == http.MethodGet:
			writeAPIResponse(w, 0, "success", &GameListResp{
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
				TotalCount: 2,
				Data: []GameInfoResp{
					{URI: "game1", Name: "Game One"},
					{URI: "game2", Name: "Game Two"},
				},
			})
		case r.URL.Path == "/developer/v1/game/create-with-version":
			var req GameCreateWithVersionReq
			json.NewDecoder(r.Body).Decode(&req)
			writeAPIResponse(w, 0, "success", &GameInfoResp{
				URI:       "newgame",
				Name:      req.Name,
				AccessUrl: "https://game.example.com",
			})
		case r.URL.Path == "/developer/v1/game/update-with-version":
			writeAPIResponse(w, 0, "success", &GameInfoResp{
				URI:       "existing",
				Name:      "Updated",
				AccessUrl: "https://game.example.com",
			})
		case strings.HasPrefix(r.URL.Path, "/developer/v1/game/") && r.Method == http.MethodPut:
			writeAPIResponse(w, 0, "success", &GameInfoResp{
				URI:       "existing",
				Name:      "Updated Info",
				AccessUrl: "https://game.example.com",
			})
		case r.URL.Path == "/developer/v1/game/apply-review":
			writeAPIResponse(w, 0, "success", nil)
		default:
			writeAPIResponse(w, 404, "not found", nil)
		}
	}))

	client := NewClient(apiServer.Client(), apiServer.URL)
	uploader := NewOSSUploader(apiServer.Client())

	return &testEnv{
		apiServer: apiServer,
		ossServer: ossServer,
		prov:      NewProvider(client, uploader),
		tmpDir:    t.TempDir(),
	}
}

func (e *testEnv) close() {
	e.apiServer.Close()
	e.ossServer.Close()
}

func (e *testEnv) writeTmpFile(name, content string) string {
	p := filepath.Join(e.tmpDir, name)
	os.WriteFile(p, []byte(content), 0644)
	return p
}

func TestProvider_CreateMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	pkgPath := env.writeTmpFile("game.zip", "fake zip")

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:    "test@example.com",
		Password:      "testpass",
		Mode:          provider.PublishModeCreate,
		GameName:      "Test Game",
		GameDescription: "A test",
		PackagePath:   pkgPath,
		Version:       "1.0.0",
		ChangeLog:     "Initial",
		ScreenConfig:  &provider.ScreenConfig{ScreenType: 1, HalfSupport: 2, HalfRatio: "0.75"},
		BuildConfig:   &provider.BuildConfig{
			Backend:  provider.BuildConfigEntry{WorkDir: "admin", Cmd: "./server"},
			Frontend: provider.BuildConfigEntry{WorkDir: "h5"},
			Socket:   provider.BuildConfigEntry{WorkDir: "logic", Cmd: "./server -type logic"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Provider != "3os" {
		t.Errorf("expected provider=3os, got %s", result.Provider)
	}
	if result.Mode != provider.PublishModeCreate {
		t.Errorf("expected mode=create, got %s", result.Mode)
	}
	if result.GameURI != "newgame" {
		t.Errorf("expected gameURI=newgame, got %s", result.GameURI)
	}
	if len(result.Assets) != 1 {
		t.Fatalf("expected 1 asset (package), got %d", len(result.Assets))
	}
	if !strings.HasSuffix(result.Assets[0].URL, ".zip") {
		t.Errorf("expected package URL ending in .zip, got %s", result.Assets[0].URL)
	}
	if result.Assets[0].Label != "package" {
		t.Errorf("expected label=package, got %s", result.Assets[0].Label)
	}
}

func TestProvider_CreateMode_WithSQL(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	pkgPath := env.writeTmpFile("game.zip", "zip")
	sqlPath := env.writeTmpFile("init.sql", "CREATE TABLE t;")

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:  "test@example.com",
		Password:    "testpass",
		Mode:        provider.PublishModeCreate,
		GameName:    "Test",
		PackagePath: pkgPath,
		SQLPath:     sqlPath,
		Version:     "1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(result.Assets))
	}
}

func TestProvider_ListMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier: "test@example.com",
		Password:   "testpass",
		Mode:       provider.PublishModeList,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != provider.PublishModeList {
		t.Errorf("expected mode=list, got %s", result.Mode)
	}
	if result.GameList == nil {
		t.Fatal("expected game list in result")
	}
	if result.GameList.TotalCount != 2 {
		t.Errorf("expected totalCount=2, got %d", result.GameList.TotalCount)
	}
	if len(result.GameList.Games) != 2 {
		t.Fatalf("expected 2 games, got %d", len(result.GameList.Games))
	}
	if result.GameList.Games[0].URI != "game1" {
		t.Errorf("expected first game uri=game1, got %s", result.GameList.Games[0].URI)
	}
}

func TestProvider_UpdateVersionMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	pkgPath := env.writeTmpFile("game-v2.zip", "zip v2")

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:  "test@example.com",
		Password:    "testpass",
		Mode:        provider.PublishModeUpdateVersion,
		GameURI:     "existing",
		PackagePath: pkgPath,
		Version:     "1.1.0",
		ChangeLog:   "Bug fixes",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != provider.PublishModeUpdateVersion {
		t.Errorf("expected mode=update-version, got %s", result.Mode)
	}
	if result.GameURI != "existing" {
		t.Errorf("expected gameURI=existing, got %s", result.GameURI)
	}
}

func TestProvider_UpdateInfoMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:    "test@example.com",
		Password:      "testpass",
		Mode:          provider.PublishModeUpdateInfo,
		GameURI:       "existing",
		GameDescription: "Updated desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != provider.PublishModeUpdateInfo {
		t.Errorf("expected mode=update-info, got %s", result.Mode)
	}
	if result.GameURI != "existing" {
		t.Errorf("expected gameURI=existing, got %s", result.GameURI)
	}
}

func TestProvider_ApplyReviewMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier: "test@example.com",
		Password:   "testpass",
		Mode:       provider.PublishModeApplyReview,
		ReviewURI:  "review-123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ReviewURI != "review-123" {
		t.Errorf("expected reviewURI=review-123, got %s", result.ReviewURI)
	}
	if !result.ReviewApplied {
		t.Error("expected reviewApplied=true")
	}
}

func TestProvider_CreateMode_WithReview(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	pkgPath := env.writeTmpFile("game.zip", "zip")

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:  "test@example.com",
		Password:    "testpass",
		Mode:        provider.PublishModeCreate,
		GameName:    "Test",
		PackagePath: pkgPath,
		Version:     "1.0.0",
		ReviewURI:   "review-auto",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.ReviewApplied {
		t.Error("expected reviewApplied=true")
	}
	if result.ReviewURI != "review-auto" {
		t.Errorf("expected reviewURI=review-auto, got %s", result.ReviewURI)
	}
}

func TestProvider_AuthFailure(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAPIResponse(w, 100, "invalid credentials", nil)
	}))
	defer apiServer.Close()

	client := NewClient(apiServer.Client(), apiServer.URL)
	uploader := NewOSSUploader(apiServer.Client())
	prov := NewProvider(client, uploader)

	_, err := prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier: "bad@example.com",
		Password:   "wrong",
		Mode:       provider.PublishModeList,
	})
	if err == nil {
		t.Fatal("expected error for auth failure")
	}
	if !strings.Contains(err.Error(), "auth failed") {
		t.Errorf("expected auth failure error, got: %v", err)
	}
}

func TestProvider_UploadFailure_StopsBeforePublish(t *testing.T) {
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/common/v1/auth/login":
			writeAPIResponse(w, 0, "success", &AuthLoginResp{AccessToken: "token"})
		case "/developer/v1/file/policy-token":
			writeAPIResponse(w, 0, "success", &FilePolicyTokenResp{
				Policy:    "p",
				Signature: "s",
				Host:      "http://localhost:0", // will fail
				Dir:       "dir/",
			})
		case "/developer/v1/game/create-with-version":
			t.Error("should not reach create-with-version after upload failure")
		}
	}))
	defer apiServer.Close()

	client := NewClient(apiServer.Client(), apiServer.URL)
	uploader := NewOSSUploader(apiServer.Client())
	prov := NewProvider(client, uploader)

	tmpFile := filepath.Join(t.TempDir(), "game.zip")
	os.WriteFile(tmpFile, []byte("content"), 0644)

	_, err := prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:  "test@example.com",
		Password:    "testpass",
		Mode:        provider.PublishModeCreate,
		GameName:    "Test",
		PackagePath: tmpFile,
		Version:     "1.0.0",
	})
	if err == nil {
		t.Fatal("expected error for upload failure")
	}
}

func TestProvider_PublishSucceeds_ReviewFails(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	// Override apply-review to fail
	env.apiServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/common/v1/auth/login":
			writeAPIResponse(w, 0, "success", &AuthLoginResp{AccessToken: "token"})
		case r.URL.Path == "/developer/v1/file/policy-token":
			writeAPIResponse(w, 0, "success", &FilePolicyTokenResp{
				Policy: "p", Signature: "s",
				Host: env.ossServer.URL, Dir: "dir/",
			})
		case r.URL.Path == "/developer/v1/game/create-with-version":
			writeAPIResponse(w, 0, "success", &GameInfoResp{
				URI: "newgame", AccessUrl: "https://game.example.com",
			})
		case r.URL.Path == "/developer/v1/game/apply-review":
			writeAPIResponse(w, 100, "review state invalid", nil)
		}
	})

	pkgPath := env.writeTmpFile("game.zip", "zip")

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier:  "test@example.com",
		Password:    "testpass",
		Mode:        provider.PublishModeCreate,
		GameName:    "Test",
		PackagePath: pkgPath,
		Version:     "1.0.0",
		ReviewURI:   "review-123",
	})
	if err == nil {
		t.Fatal("expected partial failure error")
	}
	if result == nil {
		t.Fatal("expected non-nil result for partial success")
	}
	if result.GameURI != "newgame" {
		t.Errorf("expected game info preserved in partial result, got URI=%s", result.GameURI)
	}
	if result.ReviewApplied {
		t.Error("expected reviewApplied=false")
	}
	if !strings.Contains(err.Error(), "review") {
		t.Errorf("expected review failure in error, got: %v", err)
	}
}

func TestProvider_UnsupportedMode(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	_, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier: "test@example.com",
		Password:   "testpass",
		Mode:       "invalid",
	})
	if err == nil {
		t.Fatal("expected error for unsupported mode")
	}
}

func TestProvider_ListMode_DefaultPagination(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()

	result, err := env.prov.Deploy(context.Background(), provider.DeployConfig{
		Identifier: "test@example.com",
		Password:   "testpass",
		Mode:       provider.PublishModeList,
		// Page and PageSize left at 0 — should default to 1 and 10
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GameList.Page != 1 {
		t.Errorf("expected default page=1, got %d", result.GameList.Page)
	}
	if result.GameList.PageSize != 10 {
		t.Errorf("expected default pageSize=10, got %d", result.GameList.PageSize)
	}
}

func TestProvider_Name(t *testing.T) {
	env := newTestEnv(t)
	defer env.close()
	if env.prov.Name() != "3os" {
		t.Errorf("expected provider name=3os, got %s", env.prov.Name())
	}
}
