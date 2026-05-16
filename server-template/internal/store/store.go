package store

import (
	"sync"
	"time"
)

// Store is an in-memory persistence layer for local development.
type Store struct {
	mu             sync.RWMutex
	sessions       map[string]*SessionRecord
	profiles       map[string]*ProfileRecord
	balances       map[string]*BalanceRecord
	spins          map[string][]*SpinRecord
	leaderboard    []*SlotLeaderboardEntry
}

type SessionRecord struct {
	SessionID string
	Token     string
	PlayerID  string
	ExpiresAt time.Time
}

type ProfileRecord struct {
	PlayerID  string
	Nickname  string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BalanceRecord struct {
	PlayerID  string
	Balance   int64
	UpdatedAt time.Time
}

type SpinRecord struct {
	SpinID      string
	PlayerID    string
	Wager       int64
	TotalPayout int64
	Balance     int64
	Reels       [][]string
	PaylineWins []PaylineWinRecord
	SpunAt      time.Time
}

type PaylineWinRecord struct {
	PaylineID int
	Symbol    string
	Count     int
	Payout    int64
}

type SlotLeaderboardEntry struct {
	Rank      int
	PlayerID  string
	Nickname  string
	Balance   int64
	UpdatedAt time.Time
}

func New() *Store {
	return &Store{
		sessions:    make(map[string]*SessionRecord),
		profiles:    make(map[string]*ProfileRecord),
		balances:    make(map[string]*BalanceRecord),
		spins:       make(map[string][]*SpinRecord),
		leaderboard: make([]*SlotLeaderboardEntry, 0),
	}
}

func (s *Store) SaveSession(rec *SessionRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[rec.Token] = rec
}

func (s *Store) GetSessionByToken(token string) (*SessionRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.sessions[token]
	if !ok {
		return nil, false
	}
	if time.Now().After(rec.ExpiresAt) {
		return nil, false
	}
	return rec, true
}

func (s *Store) SaveProfile(rec *ProfileRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[rec.PlayerID] = rec
}

func (s *Store) GetProfile(playerID string) (*ProfileRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.profiles[playerID]
	return rec, ok
}

func (s *Store) GetProfileByPlayerID(playerID string) (*ProfileRecord, bool) {
	return s.GetProfile(playerID)
}

func (s *Store) GetBalance(playerID string) (*BalanceRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.balances[playerID]
	return rec, ok
}

func (s *Store) InitBalance(playerID string, amount int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.balances[playerID]; !exists {
		s.balances[playerID] = &BalanceRecord{
			PlayerID:  playerID,
			Balance:   amount,
			UpdatedAt: time.Now(),
		}
	}
}

// UpdateBalance deducts wager, adds payout, and saves the spin record atomically.
func (s *Store) UpdateBalance(playerID string, wager, payout int64, rec *SpinRecord) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bal, exists := s.balances[playerID]
	if !exists {
		bal = &BalanceRecord{PlayerID: playerID, Balance: 0, UpdatedAt: time.Now()}
		s.balances[playerID] = bal
	}

	if bal.Balance < wager {
		return bal.Balance, ErrInsufficientBalance
	}

	bal.Balance = bal.Balance - wager + payout
	bal.UpdatedAt = time.Now()
	rec.Balance = bal.Balance

	s.spins[playerID] = append(s.spins[playerID], rec)
	s.rebuildLeaderboard()

	return bal.Balance, nil
}

var ErrInsufficientBalance = func() *insufficientBalanceErr {
	return &insufficientBalanceErr{}
}()

type insufficientBalanceErr struct{}

func (e *insufficientBalanceErr) Error() string { return "insufficient virtual credits" }

func (s *Store) GetSpinHistory(playerID string, limit, offset int) ([]*SpinRecord, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := s.spins[playerID]
	total := len(history)

	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	result := make([]*SpinRecord, end-start)
	copy(result, history[start:end])
	return result, total
}

func (s *Store) GetSlotLeaderboard(limit, offset int) ([]*SlotLeaderboardEntry, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.leaderboard)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	result := make([]*SlotLeaderboardEntry, end-start)
	copy(result, s.leaderboard[start:end])
	return result, total
}

func (s *Store) rebuildLeaderboard() {
	entries := make([]*SlotLeaderboardEntry, 0, len(s.balances))
	for _, br := range s.balances {
		profile, hasProfile := s.profiles[br.PlayerID]
		nickname := br.PlayerID
		if hasProfile && profile.Nickname != "" {
			nickname = profile.Nickname
		}
		entries = append(entries, &SlotLeaderboardEntry{
			PlayerID:  br.PlayerID,
			Nickname:  nickname,
			Balance:   br.Balance,
			UpdatedAt: br.UpdatedAt,
		})
	}

	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Balance > entries[i].Balance {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	for i, e := range entries {
		e.Rank = i + 1
	}

	s.leaderboard = entries
}
