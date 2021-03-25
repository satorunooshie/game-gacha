package service

import (
	"game-gacha/pkg/server/model/mock_model"
	"github.com/golang/mock/gomock"
)

type mockRepository struct {
	userRepository               *mock_model.MockUserRepositoryInterface
	collectionItemRepository     *mock_model.MockCollectionItemRepositoryInterface
	userCollectionItemRepository *mock_model.MockUserCollectionItemRepositoryInterface
}

func newMockRepository(ctrl *gomock.Controller) *mockRepository {
	return &mockRepository{
		userRepository:               mock_model.NewMockUserRepositoryInterface(ctrl),
		collectionItemRepository:     mock_model.NewMockCollectionItemRepositoryInterface(ctrl),
		userCollectionItemRepository: mock_model.NewMockUserCollectionItemRepositoryInterface(ctrl),
	}
}
