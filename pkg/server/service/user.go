package service

import (
	"fmt"

	"game-gacha/pkg/derror"
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
type userService struct {
	UserRepository model.UserRepositoryInterface
}

var _ UserServiceInterface = (*userService)(nil)

func NewUserService(userRepository model.UserRepositoryInterface) *userService {
	return &userService{
		UserRepository: userRepository,
	}
}

// TODO: make interface smaller
type UserServiceInterface interface {
	UserCreate(name string) (*userCreateResponse, error)
	UserGet(userID string) (*userGetResponse, error)
	UserUpdate(userID, name string) error
}

func (s *userService) UserCreate(name string) (*userCreateResponse, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	authToken, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	err = s.UserRepository.InsertUser(&model.User{
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

func (s *userService) UserGet(userID string) (*userGetResponse, error) {
	user, err := s.UserRepository.SelectUserByPK(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("%w. userID=%s", derror.ErrUserNotFound, userID)
	}
	return &userGetResponse{
		ID:        user.ID,
		authToken: user.AuthToken,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	}, nil
}

func (s *userService) UserUpdate(userID, name string) error {
	user, err := s.UserRepository.SelectUserByPK(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("%w. userID=%s", derror.ErrUserNotFound, userID)
	}
	user.Name = name
	if err = s.UserRepository.UpdateUserByPK(user); err != nil {
		return err
	}
	return nil
}
