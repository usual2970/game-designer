package balance

import (
	"testing"

	"github.com/example/game-designer-server/internal/store"
)

func TestGet_NoBalance(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	resp := svc.Get("unknown-player")
	if resp.Balance != 0 {
		t.Errorf("expected 0 balance for unknown player, got %d", resp.Balance)
	}
}

func TestInitAndGet(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	svc.Init("player1")
	resp := svc.Get("player1")
	if resp.Balance != DefaultBalance {
		t.Errorf("expected default balance %d, got %d", DefaultBalance, resp.Balance)
	}
}

func TestInit_Idempotent(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	svc.Init("player1")
	svc.Init("player1")
	resp := svc.Get("player1")
	if resp.Balance != DefaultBalance {
		t.Errorf("expected balance unchanged after double init, got %d", resp.Balance)
	}
}
