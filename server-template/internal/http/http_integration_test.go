package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/game-designer-server/internal/balance"
	"github.com/example/game-designer-server/internal/profile"
	"github.com/example/game-designer-server/internal/session"
	"github.com/example/game-designer-server/internal/slot"
	"github.com/example/game-designer-server/internal/store"
)

type testRNG struct {
	values []int
	index  int
}

func (t *testRNG) Intn(n int) int {
	v := t.values[t.index%len(t.values)]
	t.index++
	return v % n
}

func setupTestHandler() *Handler {
	s := store.New()
	balSvc := balance.NewService(s)
	sessSvc := session.NewService(s, time.Hour, balSvc)
	profSvc := profile.NewService(s)
	slotSvc := slot.NewService(s, balSvc)
	return NewHandler(sessSvc, profSvc, slotSvc, balSvc)
}

func setupTestHandlerWithRNG(rng slot.RNG) (*Handler, *slot.Service) {
	s := store.New()
	balSvc := balance.NewService(s)
	sessSvc := session.NewService(s, time.Hour, balSvc)
	profSvc := profile.NewService(s)
	slotSvc := slot.NewService(s, balSvc)
	slotSvc.SetRNG(rng)
	return NewHandler(sessSvc, profSvc, slotSvc, balSvc), slotSvc
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

func TestFullSlotLoop(t *testing.T) {
	// Deterministic reels: all "Cherry" (index 0)
	rng := &testRNG{values: []int{0, 0, 0, 0, 0, 0, 0, 0, 0}}
	h, _ := setupTestHandlerWithRNG(rng)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Step 1: Create session
	token := createTestSession(h, "slot-player")
	if token == "" {
		t.Fatal("expected non-empty session token")
	}

	// Step 2: Get profile
	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("X-Session-Token", token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET profile: expected 200, got %d", w.Code)
	}

	// Step 3: Get slot config
	req = httptest.NewRequest("GET", "/api/v1/slot/config", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET slot/config: expected 200, got %d", w.Code)
	}
	var configResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&configResp)
	if configResp["reels"] != float64(3) {
		t.Errorf("expected 3 reels, got %v", configResp["reels"])
	}

	// Step 4: Get balance
	req = httptest.NewRequest("GET", "/api/v1/balance", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET balance: expected 200, got %d", w.Code)
	}
	var balResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&balResp)
	if balResp["balance"] != float64(1000) {
		t.Errorf("expected initial balance=1000, got %v", balResp["balance"])
	}

	// Step 5: Spin
	spinBody, _ := json.Marshal(map[string]interface{}{"wager": float64(10)})
	req = httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(spinBody))
	req.Header.Set("X-Session-Token", token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("POST spin: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	var spinResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&spinResp)
	if spinResp["spinId"] == nil || spinResp["spinId"] == "" {
		t.Error("expected non-empty spinId")
	}
	if spinResp["totalPayout"] == nil {
		t.Error("expected totalPayout field")
	}
	if spinResp["balance"] == nil {
		t.Error("expected balance field")
	}

	// Step 6: Get spin history
	req = httptest.NewRequest("GET", "/api/v1/spin/history", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET spin/history: expected 200, got %d", w.Code)
	}
	var histResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&histResp)
	entries := histResp["entries"].([]interface{})
	if len(entries) != 1 {
		t.Errorf("expected 1 spin history entry, got %d", len(entries))
	}

	// Step 7: Get leaderboard
	req = httptest.NewRequest("GET", "/api/v1/leaderboard", nil)
	req.Header.Set("X-Session-Token", token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GET leaderboard: expected 200, got %d", w.Code)
	}
	var lbResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&lbResp)
	lbEntries := lbResp["entries"].([]interface{})
	if len(lbEntries) != 1 {
		t.Fatalf("expected 1 leaderboard entry, got %d", len(lbEntries))
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

func TestSpin_InvalidWager(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	token := createTestSession(h, "wager-player")

	body, _ := json.Marshal(map[string]interface{}{"wager": float64(0)})
	req := httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(body))
	req.Header.Set("X-Session-Token", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400 for zero wager, got %d", w.Code)
	}
	var errResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["code"] != "INVALID_PARAMETERS" {
		t.Errorf("expected code=INVALID_PARAMETERS, got %v", errResp["code"])
	}
}

func TestSpin_InsufficientBalance(t *testing.T) {
	// Use deterministic RNG that produces non-matching symbols (no wins)
	rng := &testRNG{values: []int{0, 1, 2, 3, 4, 5, 0, 1, 2}}
	h, _ := setupTestHandlerWithRNG(rng)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	token := createTestSession(h, "poor-player")

	// Drain balance with max wagers (10 * 100 = 1000, starting balance is 1000)
	for i := 0; i < 10; i++ {
		body, _ := json.Marshal(map[string]interface{}{"wager": float64(100)})
		req := httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(body))
		req.Header.Set("X-Session-Token", token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}

	// Now try one more spin with insufficient balance
	body, _ := json.Marshal(map[string]interface{}{"wager": float64(1)})
	req := httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(body))
	req.Header.Set("X-Session-Token", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400 for insufficient balance, got %d", w.Code)
	}
	var errResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp["code"] != "INSUFFICIENT_BALANCE" {
		t.Errorf("expected code=INSUFFICIENT_BALANCE, got %v", errResp["code"])
	}
}

func TestMultiplePlayers_LeaderboardRanking(t *testing.T) {
	h := setupTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	token1 := createTestSession(h, "lb-p1")
	token2 := createTestSession(h, "lb-p2")

	// Each player spins once
	body, _ := json.Marshal(map[string]interface{}{"wager": float64(10)})
	req := httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(body))
	req.Header.Set("X-Session-Token", token1)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body2, _ := json.Marshal(map[string]interface{}{"wager": float64(10)})
	req2 := httptest.NewRequest("POST", "/api/v1/spin", bytes.NewReader(body2))
	req2.Header.Set("X-Session-Token", token2)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	// Check leaderboard
	req3 := httptest.NewRequest("GET", "/api/v1/leaderboard", nil)
	req3.Header.Set("X-Session-Token", token1)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)

	var lbResp map[string]interface{}
	json.NewDecoder(w3.Body).Decode(&lbResp)
	entries := lbResp["entries"].([]interface{})
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
