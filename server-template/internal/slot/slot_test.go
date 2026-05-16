package slot

import (
	"testing"

	"github.com/example/game-designer-server/internal/balance"
	"github.com/example/game-designer-server/internal/store"
)

// deterministicRNG cycles through predefined values for reproducible tests.
type deterministicRNG struct {
	values []int
	index  int
}

func (d *deterministicRNG) Intn(n int) int {
	v := d.values[d.index%len(d.values)]
	d.index++
	return v % n
}

func setupTestSlotService(rng RNG) (*Service, *store.Store) {
	s := store.New()
	balSvc := balance.NewService(s)
	svc := NewService(s, balSvc)
	if rng != nil {
		svc.SetRNG(rng)
	}
	return svc, s
}

func TestGetConfig(t *testing.T) {
	svc, _ := setupTestSlotService(nil)
	config := svc.GetConfig()

	if config.Reels != 3 {
		t.Errorf("expected 3 reels, got %d", config.Reels)
	}
	if config.Rows != 3 {
		t.Errorf("expected 3 rows, got %d", config.Rows)
	}
	if len(config.Paylines) != 5 {
		t.Errorf("expected 5 paylines, got %d", len(config.Paylines))
	}
	if len(config.Symbols) != 7 {
		t.Errorf("expected 7 symbols, got %d", len(config.Symbols))
	}
	if config.MinWager != 1 {
		t.Errorf("expected minWager=1, got %d", config.MinWager)
	}
	if config.MaxWager != 100 {
		t.Errorf("expected maxWager=100, got %d", config.MaxWager)
	}
}

func TestSpin_HappyPath(t *testing.T) {
	// All reels show "Seven" → payline 1 (top row) wins
	// Symbol index for Seven is 5
	rng := &deterministicRNG{values: []int{5, 5, 5, 0, 0, 0, 0, 0, 0}}
	svc, s := setupTestSlotService(rng)

	s.InitBalance("player1", 1000)

	result, err := svc.Spin("player1", SpinRequest{Wager: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Wager != 10 {
		t.Errorf("expected wager=10, got %d", result.Wager)
	}
	if result.TotalPayout == 0 {
		t.Error("expected non-zero payout for winning reels")
	}
	if result.SpinID == "" {
		t.Error("expected non-empty spinId")
	}
	// Balance = 1000 - 10 + payout
	if result.Balance != 1000-10+result.TotalPayout {
		t.Errorf("balance mismatch: expected %d, got %d", 1000-10+result.TotalPayout, result.Balance)
	}
	// All reels should show "Seven"
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if row == 0 {
				if result.Reels[row][col] != "Seven" {
					t.Errorf("expected Seven at [%d][%d], got %s", row, col, result.Reels[row][col])
				}
			}
		}
	}
}

func TestSpin_NoWin(t *testing.T) {
	// Mix symbols so no payline has all matching
	rng := &deterministicRNG{values: []int{0, 1, 2, 3, 4, 5, 0, 1, 2}}
	svc, s := setupTestSlotService(rng)

	s.InitBalance("player1", 1000)

	result, err := svc.Spin("player1", SpinRequest{Wager: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPayout != 0 {
		t.Errorf("expected zero payout for no-win spin, got %d", result.TotalPayout)
	}
	if len(result.PaylineWins) != 0 {
		t.Errorf("expected no payline wins, got %d", len(result.PaylineWins))
	}
	if result.Balance != 990 {
		t.Errorf("expected balance=990 (1000-10), got %d", result.Balance)
	}
}

func TestSpin_InvalidWager_Zero(t *testing.T) {
	svc, _ := setupTestSlotService(nil)
	_, err := svc.Spin("player1", SpinRequest{Wager: 0})
	if err != ErrInvalidWager {
		t.Errorf("expected ErrInvalidWager, got %v", err)
	}
}

func TestSpin_InvalidWager_Negative(t *testing.T) {
	svc, _ := setupTestSlotService(nil)
	_, err := svc.Spin("player1", SpinRequest{Wager: -5})
	if err != ErrInvalidWager {
		t.Errorf("expected ErrInvalidWager, got %v", err)
	}
}

func TestSpin_InvalidWager_OverMax(t *testing.T) {
	svc, _ := setupTestSlotService(nil)
	_, err := svc.Spin("player1", SpinRequest{Wager: 101})
	if err != ErrInvalidWager {
		t.Errorf("expected ErrInvalidWager, got %v", err)
	}
}

func TestSpin_InsufficientBalance(t *testing.T) {
	svc, s := setupTestSlotService(nil)
	s.InitBalance("player1", 5)

	_, err := svc.Spin("player1", SpinRequest{Wager: 10})
	if err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestSpin_BalanceUnchangedOnInsufficient(t *testing.T) {
	svc, s := setupTestSlotService(nil)
	s.InitBalance("player1", 5)

	svc.Spin("player1", SpinRequest{Wager: 10})

	bal, _ := s.GetBalance("player1")
	if bal.Balance != 5 {
		t.Errorf("expected balance unchanged at 5, got %d", bal.Balance)
	}
}

func TestSpin_RepeatedSpinsAppendHistory(t *testing.T) {
	rng := &deterministicRNG{values: []int{0, 0, 0, 1, 1, 1, 2, 2, 2}}
	svc, s := setupTestSlotService(rng)

	s.InitBalance("player1", 10000)

	svc.Spin("player1", SpinRequest{Wager: 10})
	svc.Spin("player1", SpinRequest{Wager: 10})

	history := svc.GetHistory("player1", 100, 0)
	if history.Total != 2 {
		t.Errorf("expected 2 spin history records, got %d", history.Total)
	}
	if len(history.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(history.Entries))
	}
}

func TestSpin_LeaderboardUpdated(t *testing.T) {
	svc, s := setupTestSlotService(nil)

	s.InitBalance("p1", 1000)
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p1", Nickname: "Player1"})
	s.InitBalance("p2", 2000)
	s.SaveProfile(&store.ProfileRecord{PlayerID: "p2", Nickname: "Player2"})

	svc.Spin("p1", SpinRequest{Wager: 1})
	svc.Spin("p2", SpinRequest{Wager: 1})

	lb := svc.GetLeaderboard(10, 0)
	if lb.Total != 2 {
		t.Errorf("expected 2 leaderboard entries, got %d", lb.Total)
	}
	if lb.Entries[0].Rank != 1 {
		t.Errorf("expected rank 1 for first entry, got %d", lb.Entries[0].Rank)
	}
}

func TestGetHistory_Pagination(t *testing.T) {
	svc, s := setupTestSlotService(nil)
	s.InitBalance("player1", 10000)

	for i := 0; i < 5; i++ {
		svc.Spin("player1", SpinRequest{Wager: 1})
	}

	page1 := svc.GetHistory("player1", 2, 0)
	if len(page1.Entries) != 2 {
		t.Errorf("expected 2 entries on page 1, got %d", len(page1.Entries))
	}
	if page1.Total != 5 {
		t.Errorf("expected total=5, got %d", page1.Total)
	}

	page2 := svc.GetHistory("player1", 2, 2)
	if len(page2.Entries) != 2 {
		t.Errorf("expected 2 entries on page 2, got %d", len(page2.Entries))
	}
}

func TestPaylineEvaluation_MultipleWins(t *testing.T) {
	// All positions are "Cherry" (index 0) — all 5 paylines should win
	rng := &deterministicRNG{values: []int{0, 0, 0, 0, 0, 0, 0, 0, 0}}
	svc, s := setupTestSlotService(rng)

	s.InitBalance("player1", 1000)

	result, err := svc.Spin("player1", SpinRequest{Wager: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.PaylineWins) != 5 {
		t.Errorf("expected 5 payline wins for all-Cherry grid, got %d", len(result.PaylineWins))
	}
	// Each Cherry pays 5x wager = 50 per payline, 5 paylines = 250
	expectedPayout := int64(5*5) * 10
	if result.TotalPayout != expectedPayout {
		t.Errorf("expected totalPayout=%d, got %d", expectedPayout, result.TotalPayout)
	}
}
