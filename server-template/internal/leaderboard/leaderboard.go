package leaderboard

import (
	"time"

	"github.com/example/game-designer-server/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

type SubmitScoreRequest struct {
	Score    int64                  `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type SubmitScoreResponse struct {
	Accepted bool  `json:"accepted"`
	Rank     int   `json:"rank"`
	BestScore int64 `json:"bestScore"`
	IsNewBest bool  `json:"isNewBest"`
}

type LeaderboardEntry struct {
	Rank       int       `json:"rank"`
	PlayerID   string    `json:"playerId"`
	Nickname   string    `json:"nickname"`
	Score      int64     `json:"score"`
	AchievedAt time.Time `json:"achievedAt"`
}

type LeaderboardResponse struct {
	Entries []LeaderboardEntry `json:"entries"`
	Total   int                `json:"total"`
}

func (svc *Service) SubmitScore(playerID string, req SubmitScoreRequest) (*SubmitScoreResponse, error) {
	scoreRecord, isNewBest, err := svc.store.SubmitScore(playerID, req.Score)
	if err != nil {
		return nil, err
	}

	rank := 0
	entries, _ := svc.store.GetLeaderboard(1000, 0)
	for _, e := range entries {
		if e.PlayerID == playerID {
			rank = e.Rank
			break
		}
	}

	return &SubmitScoreResponse{
		Accepted:  true,
		Rank:      rank,
		BestScore: scoreRecord.BestScore,
		IsNewBest: isNewBest,
	}, nil
}

func (svc *Service) GetLeaderboard(limit, offset int) (*LeaderboardResponse, error) {
	entries, total := svc.store.GetLeaderboard(limit, offset)

	result := make([]LeaderboardEntry, len(entries))
	for i, e := range entries {
		result[i] = LeaderboardEntry{
			Rank:       e.Rank,
			PlayerID:   e.PlayerID,
			Nickname:   e.Nickname,
			Score:      e.Score,
			AchievedAt: e.AchievedAt,
		}
	}

	return &LeaderboardResponse{
		Entries: result,
		Total:   total,
	}, nil
}
