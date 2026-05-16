package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/game-designer-server/internal/gamestate"
	"github.com/example/game-designer-server/internal/leaderboard"
	"github.com/example/game-designer-server/internal/profile"
	"github.com/example/game-designer-server/internal/session"
	"github.com/example/game-designer-server/internal/store"
)

func setupTestHandler() *Handler {
	s := store.New()
	sessSvc := session.NewService(s, time.Hour)
	profSvc := profile.NewService(s)
	gsSvc := gamestate.NewService(s)
	lbSvc := leaderboard.NewService(s)
	return NewHandler(sessSvc, profSvc, gsSvc, lbSvc)
}

func createTestSession(h *Handler, playerID string) string {
	body, _ := json.Marshal(map[string]string{
		"playerId": playerID,
		"nickname": "TestPlayer",
	})
	req := httptest.NewRequest("POST", "/api/v1/session", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	return resp["token"].(string)
}

func TestFullActivityLoop(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Step 1: Create session
	token := createTestSession(h, "loop-player")
	if token == "" {
		t.Fatal("expected non-empty session token")
	}

	// Step 2: Get profile
	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("X-Session-Token", token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET profile: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// Step 3: Save game state
	saveBody, _ := json.Marshal(map[string]interface{}{
		"data":       map[string]interface{}{"level": float64(5), "coins": float64(200)},
		"checkpoint": "level-5",
	})
	req = httptest.NewRequest("PUT", "/api/v1/game-state", bytes.NewReader(saveBody))
	req.Header.Set("X-Session-Token", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("PUT game-state: expected 200, got %d", w.Code)
	}

	// Step 4: Load game state
	req = httptest.NewRequest("GET", "/api/v1/game-state", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET game-state: expected 200, got %d", w.Code)
	}
	var stateResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&stateResp)
	if stateResp["checkpoint"] != "level-5" {
		t.Errorf("expected checkpoint=level-5, got %v", stateResp["checkpoint"])
	}

	// Step 5: Submit score
	scoreBody, _ := json.Marshal(map[string]interface{}{"score": 1500})
	req = httptest.NewRequest("POST", "/api/v1/scores", bytes.NewReader(scoreBody))
	req.Header.Set("X-Session-Token", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("POST scores: expected 200, got %d", w.Code)
	}
	var scoreResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&scoreResp)
	if scoreResp["accepted"] != true {
		t.Error("expected accepted=true")
	}

	// Step 6: Get leaderboard
	req = httptest.NewRequest("GET", "/api/v1/leaderboard", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET leaderboard: expected 200, got %d", w.Code)
	}
	var lbResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&lbResp)
	entries := lbResp["entries"].([]interface{})
	if len(entries) != 1 {
		t.Fatalf("expected 1 leaderboard entry, got %d", len(entries))
	}
}

func TestMissingAuthToken(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestInvalidSessionToken(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("X-Session-Token", "invalid-token")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateSession_MissingPlayerID(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest("POST", "/api/v1/session", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var errResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["code"] != "INVALID_PARAMETERS" {
		t.Errorf("expected code=INVALID_PARAMETERS, got %v", errResp["code"])
	}
}

func TestGameState_NoContent(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	token := createTestSession(h, "empty-player")

	req := httptest.NewRequest("GET", "/api/v1/game-state", nil)
	req.Header.Set("X-Session-Token", token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for no game state, got %d", w.Code)
	}
}

func TestMultipleScores_LeaderboardRanking(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create two players
	token1 := createTestSession(h, "multi-p1")
	token2 := createTestSession(h, "multi-p2")

	// Player 1 scores 100
	body, _ := json.Marshal(map[string]interface{}{"score": 100})
	req := httptest.NewRequest("POST", "/api/v1/scores", bytes.NewReader(body))
	req.Header.Set("X-Session-Token", token1)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Player 2 scores 200
	body2, _ := json.Marshal(map[string]interface{}{"score": 200})
	req2 := httptest.NewRequest("POST", "/api/v1/scores", bytes.NewReader(body2))
	req2.Header.Set("X-Session-Token", token2)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	// Player 1 scores 300 (new best)
	body3, _ := json.Marshal(map[string]interface{}{"score": 300})
	req3 := httptest.NewRequest("POST", "/api/v1/scores", bytes.NewReader(body3))
	req3.Header.Set("X-Session-Token", token1)
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)
	io.ReadAll(w3.Body)

	// Check leaderboard
	req4 := httptest.NewRequest("GET", "/api/v1/leaderboard", nil)
	req4.Header.Set("X-Session-Token", token1)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, req4)

	var lbResp map[string]interface{}
	json.NewDecoder(w4.Body).Decode(&lbResp)
	entries := lbResp["entries"].([]interface{})

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	first := entries[0].(map[string]interface{})
	if first["score"] != float64(300) {
		t.Errorf("expected top score=300, got %v", first["score"])
	}
}
