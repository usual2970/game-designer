package gamestate

import (
	"testing"

	"github.com/example/game-designer-server/internal/store"
)

func TestSaveAndLoad(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	data := map[string]interface{}{
		"level":     float64(3),
		"coins":     float64(150),
		"inventory": []interface{}{"sword", "shield"},
	}

	_, err := svc.Save("player1", SaveRequest{
		Data:       data,
		Checkpoint: "level-3",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, exists, err := svc.Load("player1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected state to exist")
	}
	if resp.Checkpoint != "level-3" {
		t.Errorf("expected checkpoint=level-3, got %s", resp.Checkpoint)
	}
	if resp.Data["level"] != float64(3) {
		t.Errorf("expected level=3, got %v", resp.Data["level"])
	}
}

func TestLoad_NotExists(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	_, exists, err := svc.Load("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected exists=false for missing state")
	}
}

func TestSave_Overwrite(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	svc.Save("player1", SaveRequest{
		Data: map[string]interface{}{"level": float64(1)},
	})

	resp, _, _ := svc.Load("player1")
	if resp.Data["level"] != float64(1) {
		t.Fatalf("expected level=1, got %v", resp.Data["level"])
	}

	svc.Save("player1", SaveRequest{
		Data: map[string]interface{}{"level": float64(5)},
	})

	resp, _, _ = svc.Load("player1")
	if resp.Data["level"] != float64(5) {
		t.Errorf("expected level=5 after overwrite, got %v", resp.Data["level"])
	}
}
