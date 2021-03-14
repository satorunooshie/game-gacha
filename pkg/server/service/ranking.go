package service

import (
	"game-gacha/pkg/server/model"
)

type rankingListResponse struct {
	Ranks []*rank
}
type rank struct {
	UserID   string
	UserName string
	Rank     int
	Score    int
}

func RankingList(userID string, startPosition, limit int) (*rankingListResponse, error) {
	if _, err := model.SelectUserByPK(userID); err != nil {
		return nil, err
	}
	users, err := model.SelectUsersOrderByHighScore(startPosition-1, limit)
	if err != nil {
		return nil, err
	}
	ranks := make([]*rank, 0, len(users))
	for i, user := range users {
		rank := &rank{
			UserID:   user.ID,
			UserName: user.Name,
			Rank:     i + startPosition,
			Score:    user.HighScore,
		}
		ranks = append(ranks, rank)
	}
	return &rankingListResponse{
		Ranks: ranks,
	}, nil
}
