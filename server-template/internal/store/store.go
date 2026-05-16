package store

import (
	"sync"
	"time"
)

// Store is an in-memory persistence layer for local development.
type Store struct {
	mu          sync.RWMutex
	sessions    map[string]*SessionRecord
	profiles    map[string]*ProfileRecord
	gameStates  map[string]*GameStateRecord
	scores      map[string]*ScoreRecord // keyed by playerID
	leaderboard []*LeaderboardEntry
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

type GameStateRecord struct {
	PlayerID  string
	Data      map[string]interface{}
	Checkpoint string
	SavedAt   time.Time
}

type ScoreRecord struct {
	PlayerID  string
	BestScore int64
	UpdatedAt time.Time
}

type LeaderboardEntry struct {
	Rank      int
	PlayerID  string
	Nickname  string
	Score     int64
	AchievedAt time.Time
}

func New() *Store {
	return &Store{
		sessions:    make(map[string]*SessionRecord),
		profiles:    make(map[string]*ProfileRecord),
		gameStates:  make(map[string]*GameStateRecord),
		scores:      make(map[string]*ScoreRecord),
		leaderboard: make([]*LeaderboardEntry, 0),
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

func (s *Store) SaveGameState(rec *GameStateRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gameStates[rec.PlayerID] = rec
}

func (s *Store) GetGameState(playerID string) (*GameStateRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.gameStates[playerID]
	return rec, ok
}

func (s *Store) SubmitScore(playerID string, score int64) (*ScoreRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.scores[playerID]
	isNewBest := !exists || score > existing.BestScore

	bestScore := score
	if exists && !isNewBest {
		bestScore = existing.BestScore
	}

	now := time.Now()
	s.scores[playerID] = &ScoreRecord{
		PlayerID:  playerID,
		BestScore: bestScore,
		UpdatedAt: now,
	}

	s.rebuildLeaderboard()

	return s.scores[playerID], isNewBest, nil
}

func (s *Store) GetLeaderboard(limit, offset int) ([]*LeaderboardEntry, int) {
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

	result := make([]*LeaderboardEntry, end-start)
	copy(result, s.leaderboard[start:end])
	return result, total
}

func (s *Store) GetScore(playerID string) (*ScoreRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.scores[playerID]
	return rec, ok
}

func (s *Store) GetProfileByPlayerID(playerID string) (*ProfileRecord, bool) {
	return s.GetProfile(playerID)
}

func (s *Store) rebuildLeaderboard() {
	entries := make([]*LeaderboardEntry, 0, len(s.scores))
	for _, sr := range s.scores {
		profile, hasProfile := s.profiles[sr.PlayerID]
		nickname := sr.PlayerID
		if hasProfile && profile.Nickname != "" {
			nickname = profile.Nickname
		}
		entries = append(entries, &LeaderboardEntry{
			PlayerID:   sr.PlayerID,
			Nickname:   nickname,
			Score:      sr.BestScore,
			AchievedAt: sr.UpdatedAt,
		})
	}

	// Sort descending by score
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Score > entries[i].Score {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	for i, e := range entries {
		e.Rank = i + 1
	}

	s.leaderboard = entries
}
