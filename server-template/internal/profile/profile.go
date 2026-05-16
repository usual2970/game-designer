package profile

import (
	"errors"
	"time"

	"github.com/example/game-designer-server/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

type ProfileResponse struct {
	PlayerID  string    `json:"playerId"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

var (
	ErrNotFound = errors.New("profile not found")
)

func (svc *Service) Get(playerID string) (*ProfileResponse, error) {
	rec, ok := svc.store.GetProfile(playerID)
	if !ok {
		return nil, ErrNotFound
	}
	return &ProfileResponse{
		PlayerID:  rec.PlayerID,
		Nickname:  rec.Nickname,
		AvatarURL: rec.AvatarURL,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}, nil
}

func (svc *Service) Update(playerID string, req UpdateProfileRequest) (*ProfileResponse, error) {
	rec, ok := svc.store.GetProfile(playerID)
	if !ok {
		return nil, ErrNotFound
	}

	if req.Nickname != "" {
		rec.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		rec.AvatarURL = req.AvatarURL
	}
	rec.UpdatedAt = time.Now()

	svc.store.SaveProfile(rec)

	return &ProfileResponse{
		PlayerID:  rec.PlayerID,
		Nickname:  rec.Nickname,
		AvatarURL: rec.AvatarURL,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: rec.UpdatedAt,
	}, nil
}
