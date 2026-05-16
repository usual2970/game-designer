package profile

import (
	"testing"

	"github.com/example/game-designer-server/internal/store"
)

func TestGet_ProfileExists(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	s.SaveProfile(&store.ProfileRecord{
		PlayerID: "p1",
		Nickname: "Alice",
	})

	resp, err := svc.Get("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Nickname != "Alice" {
		t.Errorf("expected nickname=Alice, got %s", resp.Nickname)
	}
}

func TestGet_ProfileNotFound(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	_, err := svc.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate_Profile(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	s.SaveProfile(&store.ProfileRecord{
		PlayerID: "p1",
		Nickname: "Old",
	})

	resp, err := svc.Update("p1", UpdateProfileRequest{Nickname: "New"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Nickname != "New" {
		t.Errorf("expected nickname=New, got %s", resp.Nickname)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	_, err := svc.Update("nonexistent", UpdateProfileRequest{Nickname: "New"})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
