package service

import (
	"fmt"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/server/model"
)

type gameFinishResponse struct {
	GottenCoin int
}

func GameFinish(userID string, score int) (*gameFinishResponse, error) {
	gottenCoin := int(float64(score) * constant.CoinExchangeRateForScore)
	user, err := model.SelectUserByPK(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found. userID=%s", userID)
	}
	user.Coin += gottenCoin
	if user.HighScore < score {
		user.HighScore = score
	}
	if err = model.UpdateUserByPK(user); err != nil {
		return nil, err
	}
	return &gameFinishResponse{
		GottenCoin: gottenCoin,
	}, nil
}
