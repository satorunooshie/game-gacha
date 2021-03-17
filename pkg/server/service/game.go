package service

import (
	"fmt"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/server/model"
)

type gameFinishResponse struct {
	GottenCoin int
}
type gameService struct {
	UserRepository model.UserRepositoryInterface
}
type GameServiceInterface interface {
	GameFinish(userID string, score int) (*gameFinishResponse, error)
}

func NewGameService(userRepository model.UserRepositoryInterface) *gameService {
	return &gameService{
		UserRepository: userRepository,
	}
}

func (s *gameService) GameFinish(userID string, score int) (*gameFinishResponse, error) {
	gottenCoin := int(float64(score) * constant.CoinExchangeRateForScore)
	user, err := s.UserRepository.SelectUserByPK(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("%w. userID=%s", derror.ErrUserNotFound, userID)
	}
	user.Coin += gottenCoin
	if user.HighScore < score {
		user.HighScore = score
	}
	if err = s.UserRepository.UpdateUserByPK(user); err != nil {
		return nil, err
	}
	return &gameFinishResponse{
		GottenCoin: gottenCoin,
	}, nil
}
