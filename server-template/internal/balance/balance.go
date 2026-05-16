package balance

import (
	"github.com/example/game-designer-server/internal/slot"
	"github.com/example/game-designer-server/internal/store"
)

const DefaultBalance = slot.DefaultBalance

type Service struct {
	store           *store.Store
	defaultBalance  int64
}

func NewService(s *store.Store) *Service {
	return &Service{store: s, defaultBalance: DefaultBalance}
}

type BalanceResponse struct {
	Balance int64 `json:"balance"`
}

func (svc *Service) Init(playerID string) {
	svc.store.InitBalance(playerID, svc.defaultBalance)
}

func (svc *Service) Get(playerID string) *BalanceResponse {
	rec, ok := svc.store.GetBalance(playerID)
	if !ok {
		return &BalanceResponse{Balance: 0}
	}
	return &BalanceResponse{Balance: rec.Balance}
}
