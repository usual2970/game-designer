package session

import (
	"testing"
	"time"

	"github.com/example/game-designer-server/internal/store"
)

func TestCreateOrResume_NewPlayer(t *testing.T) {
	s := store.New()
	svc := NewService(s, time.Hour)

	resp, err := svc.CreateOrResume(CreateSessionRequest{
		PlayerID: "player1",
		Nickname: "TestPlayer",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PlayerID != "player1" {
		t.Errorf("expected playerId=player1, got %s", resp.PlayerID)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if !resp.IsNew {
		t.Error("expected isNew=true for new player")
	}
}

func TestCreateOrResume_ResumeExisting(t *testing.T) {
	s := store.New()
	svc := NewService(s, time.Hour)

	first, _ := svc.CreateOrResume(CreateSessionRequest{
		PlayerID: "player1",
		Nickname: "TestPlayer",
	})
	second, _ := svc.CreateOrResume(CreateSessionRequest{
		PlayerID: "player1",
	})

	if second.IsNew {
		t.Error("expected isNew=false for returning player")
	}
	if first.PlayerID != second.PlayerID {
		t.Error("player IDs should match")
	}
}

func TestCreateOrResume_MissingPlayerID(t *testing.T) {
	s := store.New()
	svc := NewService(s, time.Hour)

	_, err := svc.CreateOrResume(CreateSessionRequest{})
	if err != ErrMissingPlayerID {
		t.Errorf("expected ErrMissingPlayerID, got %v", err)
	}
}

func TestValidateToken_Valid(t *testing.T) {
	s := store.New()
	svc := NewService(s, time.Hour)

	resp, _ := svc.CreateOrResume(CreateSessionRequest{PlayerID: "player1"})
	playerID, ok := svc.ValidateToken(resp.Token)
	if !ok {
		t.Error("expected valid token")
	}
	if playerID != "player1" {
		t.Errorf("expected playerID=player1, got %s", playerID)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	s := store.New()
	svc := NewService(s, time.Hour)

	_, ok := svc.ValidateToken("nonexistent")
	if ok {
		t.Error("expected invalid token to fail")
	}
}
