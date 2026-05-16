package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/game-designer-server/internal/balance"
	"github.com/example/game-designer-server/internal/profile"
	"github.com/example/game-designer-server/internal/session"
	"github.com/example/game-designer-server/internal/slot"
	"github.com/example/game-designer-server/internal/store"
)

type Handler struct {
	sessions *session.Service
	profiles *profile.Service
	slot     *slot.Service
	balance  *balance.Service
}

func NewHandler(
	sessSvc *session.Service,
	profSvc *profile.Service,
	slotSvc *slot.Service,
	balSvc *balance.Service,
) *Handler {
	return &Handler{
		sessions: sessSvc,
		profiles: profSvc,
		slot:     slotSvc,
		balance:  balSvc,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/session", h.createSession)
	mux.HandleFunc("GET /api/v1/profile", h.auth(h.getProfile))
	mux.HandleFunc("PUT /api/v1/profile", h.auth(h.updateProfile))
	mux.HandleFunc("GET /api/v1/slot/config", h.auth(h.getSlotConfig))
	mux.HandleFunc("GET /api/v1/balance", h.auth(h.getBalance))
	mux.HandleFunc("POST /api/v1/spin", h.auth(h.spin))
	mux.HandleFunc("GET /api/v1/spin/history", h.auth(h.getSpinHistory))
	mux.HandleFunc("GET /api/v1/leaderboard", h.auth(h.getLeaderboard))
}

func (h *Handler) auth(next func(w http.ResponseWriter, r *http.Request, playerID string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing session token", nil)
			return
		}

		playerID, ok := h.sessions.ValidateToken(token)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired session token", nil)
			return
		}

		next(w, r, playerID)
	}
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	var req session.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", "Invalid request body", nil)
		return
	}

	resp, err := h.sessions.CreateOrResume(req)
	if err != nil {
		if errors.Is(err, session.ErrMissingPlayerID) {
			writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create session", nil)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request, playerID string) {
	resp, err := h.profiles.Get(playerID)
	if err != nil {
		if errors.Is(err, profile.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Profile not found", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get profile", nil)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request, playerID string) {
	var req profile.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", "Invalid request body", nil)
		return
	}

	resp, err := h.profiles.Update(playerID, req)
	if err != nil {
		if errors.Is(err, profile.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Profile not found", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update profile", nil)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) getSlotConfig(w http.ResponseWriter, r *http.Request, playerID string) {
	writeJSON(w, http.StatusOK, h.slot.GetConfig())
}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request, playerID string) {
	resp := h.balance.Get(playerID)
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) spin(w http.ResponseWriter, r *http.Request, playerID string) {
	var req slot.SpinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", "Invalid request body", nil)
		return
	}

	result, err := h.slot.Spin(playerID, req)
	if err != nil {
		if errors.Is(err, slot.ErrInvalidWager) {
			writeError(w, http.StatusBadRequest, "INVALID_PARAMETERS", err.Error(), nil)
			return
		}
		if errors.Is(err, store.ErrInsufficientBalance) {
			writeError(w, http.StatusBadRequest, "INSUFFICIENT_BALANCE", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to resolve spin", nil)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getSpinHistory(w http.ResponseWriter, r *http.Request, playerID string) {
	limit := intQuery(r, "limit", 20)
	offset := intQuery(r, "offset", 0)
	resp := h.slot.GetHistory(playerID, limit, offset)
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request, playerID string) {
	limit := intQuery(r, "limit", 10)
	offset := intQuery(r, "offset", 0)
	resp := h.slot.GetLeaderboard(limit, offset)
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   message,
		"code":    code,
		"details": details,
	})
}

func intQuery(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := parseSimpleInt(v)
	if err != nil {
		return defaultVal
	}
	return n
}

func parseSimpleInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
