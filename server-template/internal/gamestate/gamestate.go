package gamestate

import (
	"time"

	"github.com/example/game-designer-server/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

type SaveRequest struct {
	Data       map[string]interface{} `json:"data"`
	Checkpoint string                 `json:"checkpoint,omitempty"`
}

type StateResponse struct {
	Data       map[string]interface{} `json:"data"`
	Checkpoint string                 `json:"checkpoint,omitempty"`
	SavedAt    time.Time              `json:"savedAt"`
}

func (svc *Service) Save(playerID string, req SaveRequest) (*StateResponse, error) {
	now := time.Now()
	rec := &store.GameStateRecord{
		PlayerID:  playerID,
		Data:      req.Data,
		Checkpoint: req.Checkpoint,
		SavedAt:   now,
	}
	svc.store.SaveGameState(rec)

	return &StateResponse{
		Data:       req.Data,
		Checkpoint: req.Checkpoint,
		SavedAt:    now,
	}, nil
}

func (svc *Service) Load(playerID string) (*StateResponse, bool, error) {
	rec, ok := svc.store.GetGameState(playerID)
	if !ok {
		return nil, false, nil
	}
	return &StateResponse{
		Data:       rec.Data,
		Checkpoint: rec.Checkpoint,
		SavedAt:    rec.SavedAt,
	}, true, nil
}
