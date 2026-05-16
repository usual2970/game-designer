package session

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/example/game-designer-server/internal/store"
)

type Service struct {
	store *store.Store
	ttl   time.Duration
}

func NewService(s *store.Store, ttl time.Duration) *Service {
	return &Service{store: s, ttl: ttl}
}

type CreateSessionRequest struct {
	PlayerID  string `json:"playerId"`
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type SessionResponse struct {
	SessionID string    `json:"sessionId"`
	Token     string    `json:"token"`
	PlayerID  string    `json:"playerId"`
	IsNew     bool      `json:"isNew"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (svc *Service) CreateOrResume(req CreateSessionRequest) (*SessionResponse, error) {
	if req.PlayerID == "" {
		return nil, ErrMissingPlayerID
	}

	expiresAt := time.Now().Add(svc.ttl)
	sessionID := generateID()
	token := generateID()

	isNew := true
	if profile, exists := svc.store.GetProfile(req.PlayerID); exists {
		isNew = profile.Nickname == "" && req.Nickname == ""
		_ = profile
	}

	svc.store.SaveSession(&store.SessionRecord{
		SessionID: sessionID,
		Token:     token,
		PlayerID:  req.PlayerID,
		ExpiresAt: expiresAt,
	})

	if isNew {
		now := time.Now()
		svc.store.SaveProfile(&store.ProfileRecord{
			PlayerID:  req.PlayerID,
			Nickname:  req.Nickname,
			AvatarURL: req.AvatarURL,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	return &SessionResponse{
		SessionID: sessionID,
		Token:     token,
		PlayerID:  req.PlayerID,
		IsNew:     isNew,
		ExpiresAt: expiresAt,
	}, nil
}

func (svc *Service) ValidateToken(token string) (playerID string, ok bool) {
	rec, ok := svc.store.GetSessionByToken(token)
	if !ok {
		return "", false
	}
	return rec.PlayerID, true
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
