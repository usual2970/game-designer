package slot

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/example/game-designer-server/internal/store"
)

const (
	MinWager = 1
	MaxWager = 100
	NumReels = 3
	NumRows  = 3
)

type Symbol struct {
	Name             string
	PayoutMultiplier int64
}

var DefaultSymbols = []Symbol{
	{Name: "Cherry", PayoutMultiplier: 5},
	{Name: "Lemon", PayoutMultiplier: 3},
	{Name: "Orange", PayoutMultiplier: 8},
	{Name: "Plum", PayoutMultiplier: 4},
	{Name: "Bell", PayoutMultiplier: 15},
	{Name: "Seven", PayoutMultiplier: 50},
	{Name: "BAR", PayoutMultiplier: 25},
}

type Payline struct {
	ID        int
	Positions []int // row index per reel column
}

var DefaultPaylines = []Payline{
	{ID: 1, Positions: []int{0, 0, 0}}, // top row
	{ID: 2, Positions: []int{1, 1, 1}}, // middle row
	{ID: 3, Positions: []int{2, 2, 2}}, // bottom row
	{ID: 4, Positions: []int{0, 1, 2}}, // diagonal ↘
	{ID: 5, Positions: []int{2, 1, 0}}, // diagonal ↗
}

// RNG generates random integers in [0, n).
type RNG interface {
	Intn(n int) int
}

type cryptoRNG struct{}

func (cryptoRNG) Intn(n int) int {
	b := make([]byte, 4)
	rand.Read(b)
	v := int(uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3]))
	if n > 0 {
		return v % n
	}
	return 0
}

var (
	ErrInvalidWager     = errors.New("wager must be between 1 and 100")
	ErrInsufficientBalance = errors.New("insufficient virtual credits")
)

type Service struct {
	store   *store.Store
	balance BalanceChecker
	rng     RNG
	symbols []Symbol
	paylines []Payline
}

type BalanceChecker interface {
	Init(playerID string)
}

func NewService(s *store.Store, bc BalanceChecker) *Service {
	return &Service{
		store:    s,
		balance:  bc,
		rng:      cryptoRNG{},
		symbols:  DefaultSymbols,
		paylines: DefaultPaylines,
	}
}

func (svc *Service) SetRNG(rng RNG) {
	svc.rng = rng
}

type SlotConfigResponse struct {
	Reels          int              `json:"reels"`
	Rows           int              `json:"rows"`
	Paylines       []PaylineConfig  `json:"paylines"`
	Symbols        []SymbolConfig   `json:"symbols"`
	MinWager       int              `json:"minWager"`
	MaxWager       int              `json:"maxWager"`
	DefaultBalance int              `json:"defaultBalance"`
}

type PaylineConfig struct {
	ID        int   `json:"id"`
	Positions []int `json:"positions"`
}

type SymbolConfig struct {
	Name             string `json:"name"`
	PayoutMultiplier int64  `json:"payoutMultiplier"`
}

func (svc *Service) GetConfig() *SlotConfigResponse {
	paylines := make([]PaylineConfig, len(svc.paylines))
	for i, p := range svc.paylines {
		paylines[i] = PaylineConfig{ID: p.ID, Positions: p.Positions}
	}
	symbols := make([]SymbolConfig, len(svc.symbols))
	for i, s := range svc.symbols {
		symbols[i] = SymbolConfig{Name: s.Name, PayoutMultiplier: s.PayoutMultiplier}
	}
	return &SlotConfigResponse{
		Reels:          NumReels,
		Rows:           NumRows,
		Paylines:       paylines,
		Symbols:        symbols,
		MinWager:       MinWager,
		MaxWager:       MaxWager,
		DefaultBalance: 1000,
	}
}

type SpinRequest struct {
	Wager int64 `json:"wager"`
}

type PaylineWinResult struct {
	PaylineID int    `json:"paylineId"`
	Symbol    string `json:"symbol"`
	Count     int    `json:"count"`
	Payout    int64  `json:"payout"`
}

type SpinResult struct {
	SpinID      string            `json:"spinId"`
	Wager       int64             `json:"wager"`
	Reels       [][]string        `json:"reels"`
	PaylineWins []PaylineWinResult `json:"paylineWins"`
	TotalPayout int64             `json:"totalPayout"`
	Balance     int64             `json:"balance"`
}

func (svc *Service) Spin(playerID string, req SpinRequest) (*SpinResult, error) {
	if req.Wager < MinWager || req.Wager > MaxWager {
		return nil, ErrInvalidWager
	}

	svc.balance.Init(playerID)

	balRec, hasBal := svc.store.GetBalance(playerID)
	if hasBal && balRec.Balance < req.Wager {
		return nil, ErrInsufficientBalance
	}
	if !hasBal {
		return nil, ErrInsufficientBalance
	}

	reels := svc.generateReels()
	wins := svc.evaluatePaylines(reels, req.Wager)
	totalPayout := int64(0)
	for _, w := range wins {
		totalPayout += w.Payout
	}

	winRecords := make([]store.PaylineWinRecord, len(wins))
	for i, w := range wins {
		winRecords[i] = store.PaylineWinRecord{
			PaylineID: w.PaylineID,
			Symbol:    w.Symbol,
			Count:     w.Count,
			Payout:    w.Payout,
		}
	}

	spinRec := &store.SpinRecord{
		SpinID:      generateID(),
		PlayerID:    playerID,
		Wager:       req.Wager,
		TotalPayout: totalPayout,
		Reels:       reels,
		PaylineWins: winRecords,
		SpunAt:      time.Now(),
	}

	newBalance, err := svc.store.UpdateBalance(playerID, req.Wager, totalPayout, spinRec)
	if err != nil {
		return nil, err
	}

	return &SpinResult{
		SpinID:      spinRec.SpinID,
		Wager:       req.Wager,
		Reels:       reels,
		PaylineWins: wins,
		TotalPayout: totalPayout,
		Balance:     newBalance,
	}, nil
}

type SpinHistoryEntry struct {
	SpinID      string             `json:"spinId"`
	Wager       int64              `json:"wager"`
	TotalPayout int64              `json:"totalPayout"`
	Balance     int64              `json:"balance"`
	Reels       [][]string         `json:"reels"`
	PaylineWins []PaylineWinResult `json:"paylineWins"`
	SpunAt      time.Time          `json:"spunAt"`
}

type SpinHistoryResponse struct {
	Entries []SpinHistoryEntry `json:"entries"`
	Total   int                `json:"total"`
}

func (svc *Service) GetHistory(playerID string, limit, offset int) *SpinHistoryResponse {
	records, total := svc.store.GetSpinHistory(playerID, limit, offset)
	entries := make([]SpinHistoryEntry, len(records))
	for i, r := range records {
		wins := make([]PaylineWinResult, len(r.PaylineWins))
		for j, w := range r.PaylineWins {
			wins[j] = PaylineWinResult{
				PaylineID: w.PaylineID,
				Symbol:    w.Symbol,
				Count:     w.Count,
				Payout:    w.Payout,
			}
		}
		entries[i] = SpinHistoryEntry{
			SpinID:      r.SpinID,
			Wager:       r.Wager,
			TotalPayout: r.TotalPayout,
			Balance:     r.Balance,
			Reels:       r.Reels,
			PaylineWins: wins,
			SpunAt:      r.SpunAt,
		}
	}
	return &SpinHistoryResponse{Entries: entries, Total: total}
}

type SlotLeaderboardEntry struct {
	Rank      int       `json:"rank"`
	PlayerID  string    `json:"playerId"`
	Nickname  string    `json:"nickname"`
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SlotLeaderboardResponse struct {
	Entries []SlotLeaderboardEntry `json:"entries"`
	Total   int                    `json:"total"`
}

func (svc *Service) GetLeaderboard(limit, offset int) *SlotLeaderboardResponse {
	entries, total := svc.store.GetSlotLeaderboard(limit, offset)
	result := make([]SlotLeaderboardEntry, len(entries))
	for i, e := range entries {
		result[i] = SlotLeaderboardEntry{
			Rank:      e.Rank,
			PlayerID:  e.PlayerID,
			Nickname:  e.Nickname,
			Balance:   e.Balance,
			UpdatedAt: e.UpdatedAt,
		}
	}
	return &SlotLeaderboardResponse{Entries: result, Total: total}
}

func (svc *Service) generateReels() [][]string {
	grid := make([][]string, NumRows)
	for row := 0; row < NumRows; row++ {
		grid[row] = make([]string, NumReels)
		for col := 0; col < NumReels; col++ {
			grid[row][col] = svc.symbols[svc.rng.Intn(len(svc.symbols))].Name
		}
	}
	return grid
}

func (svc *Service) evaluatePaylines(reels [][]string, wager int64) []PaylineWinResult {
	var wins []PaylineWinResult
	for _, pl := range svc.paylines {
		firstSymbol := reels[pl.Positions[0]][0]
		allMatch := true
		for col := 1; col < len(pl.Positions); col++ {
			if reels[pl.Positions[col]][col] != firstSymbol {
				allMatch = false
				break
			}
		}
		if allMatch {
			var multiplier int64
			for _, s := range svc.symbols {
				if s.Name == firstSymbol {
					multiplier = s.PayoutMultiplier
					break
				}
			}
			wins = append(wins, PaylineWinResult{
				PaylineID: pl.ID,
				Symbol:    firstSymbol,
				Count:     len(pl.Positions),
				Payout:    multiplier * wager,
			})
		}
	}
	return wins
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
