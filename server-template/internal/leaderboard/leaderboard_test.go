package leaderboard

import (
	"testing"

	"github.com/example/game-designer-server/internal/store"
)

func TestSubmitScore_NewPlayer(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	resp, err := svc.SubmitScore("player1", SubmitScoreRequest{Score: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Accepted {
		t.Error("expected accepted=true")
	}
	if resp.BestScore != 100 {
		t.Errorf("expected bestScore=100, got %d", resp.BestScore)
	}
	if !resp.IsNewBest {
		t.Error("expected isNewBest=true for first submission")
	}
	if resp.Rank != 1 {
		t.Errorf("expected rank=1, got %d", resp.Rank)
	}
}

func TestSubmitScore_HigherScoreUpdatesBest(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	svc.SubmitScore("player1", SubmitScoreRequest{Score: 100})
	resp, _ := svc.SubmitScore("player1", SubmitScoreRequest{Score: 200})

	if resp.BestScore != 200 {
		t.Errorf("expected bestScore=200, got %d", resp.BestScore)
	}
	if !resp.IsNewBest {
		t.Error("expected isNewBest=true for higher score")
	}
}

func TestSubmitScore_LowerScoreKeepsBest(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	svc.SubmitScore("player1", SubmitScoreRequest{Score: 200})
	resp, _ := svc.SubmitScore("player1", SubmitScoreRequest{Score: 100})

	if resp.BestScore != 200 {
		t.Errorf("expected bestScore=200, got %d", resp.BestScore)
	}
	if resp.IsNewBest {
		t.Error("expected isNewBest=false for lower score")
	}
}

func TestLeaderboard_Ordering(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	// Save profiles for nickname resolution
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p1", Nickname: "Alice"})
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p2", Nickname: "Bob"})
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p3", Nickname: "Charlie"})

	svc.SubmitScore("p1", SubmitScoreRequest{Score: 300})
	svc.SubmitScore("p2", SubmitScoreRequest{Score: 100})
	svc.SubmitScore("p3", SubmitScoreRequest{Score: 200})

	resp, err := svc.GetLeaderboard(10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("expected total=3, got %d", resp.Total)
	}
	if len(resp.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(resp.Entries))
	}

	if resp.Entries[0].Nickname != "Alice" || resp.Entries[0].Score != 300 || resp.Entries[0].Rank != 1 {
		t.Errorf("rank 1 mismatch: %+v", resp.Entries[0])
	}
	if resp.Entries[1].Nickname != "Charlie" || resp.Entries[1].Score != 200 || resp.Entries[1].Rank != 2 {
		t.Errorf("rank 2 mismatch: %+v", resp.Entries[1])
	}
	if resp.Entries[2].Nickname != "Bob" || resp.Entries[2].Score != 100 || resp.Entries[2].Rank != 3 {
		t.Errorf("rank 3 mismatch: %+v", resp.Entries[2])
	}
}

func TestLeaderboard_Pagination(t *testing.T) {
	s := store.New()
	svc := NewService(s)

	s.SaveProfile(&store.ProfileRecord{PlayerID: "p1", Nickname: "A"})
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p2", Nickname: "B"})
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p3", Nickname: "C"})

	svc.SubmitScore("p1", SubmitScoreRequest{Score: 300})
	svc.SubmitScore("p2", SubmitScoreRequest{Score: 200})
	svc.SubmitScore("p3", SubmitScoreRequest{Score: 100})

	resp, _ := svc.GetLeaderboard(2, 0)
	if len(resp.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(resp.Entries))
	}
	if resp.Total != 3 {
		t.Errorf("expected total=3, got %d", resp.Total)
	}

	resp2, _ := svc.GetLeaderboard(2, 2)
	if len(resp2.Entries) != 1 {
		t.Fatalf("expected 1 entry for offset=2, got %d", len(resp2.Entries))
	}
}
