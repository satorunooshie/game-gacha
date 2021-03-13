package service

import (
	"fmt"
	"game-gacha/pkg/server/model"
	"github.com/google/uuid"
)

type userCreateResponse struct {
	Token string
}
type userGetResponse struct {
	ID        string
	Name      string
	authToken string
	HighScore int
	Coin      int
}
type userUpdateResponse struct{}

func UserCreate(name string) (*userCreateResponse, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	authToken, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	err = model.InsertUser(&model.User{
		ID:        userID.String(),
		AuthToken: authToken.String(),
		Name:      name,
		HighScore: 0,
		Coin:      0,
	})
	if err != nil {
		return nil, err
	}
	return &userCreateResponse{Token: authToken.String()}, nil
}

func UserGet(userID string) (*userGetResponse, error) {
	user, err := model.SelectUserByPK(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found. userID=%s", userID)
	}
	return &userGetResponse{
		ID:        user.ID,
		authToken: user.AuthToken,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	}, nil
}

func UserUpdate(userID string, name string) error {
	user, err := model.SelectUserByPK(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found. userID=%s", userID)
	}
	user.Name = name
	if err = model.UpdateUserByPK(user); err != nil {
		return err
	}
	return nil
}
