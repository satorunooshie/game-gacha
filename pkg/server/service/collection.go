package service

import (
	"game-gacha/pkg/server/model"
)

type collectionListResponse struct {
	Collections []*collection
}
type collection struct {
	CollectionID string
	Name         string
	Rarity       int
	HasItem      bool
}
type collectionService struct {
	UserRepository               model.UserRepositoryInterface
	UserCollectionItemRepository model.UserCollectionItemRepositoryInterface
	CollectionItemRepository     model.CollectionItemRepositoryInterface
}
type CollectionServiceInterface interface {
	CollectionList(userID string) (*collectionListResponse, error)
}

var _ CollectionServiceInterface = (*collectionService)(nil)

func NewCollectionService(
	userRepository model.UserRepositoryInterface,
	userCollectionItemRepository model.UserCollectionItemRepositoryInterface,
	collectionItemRepository model.CollectionItemRepositoryInterface,
) *collectionService {
	return &collectionService{
		UserRepository:               userRepository,
		UserCollectionItemRepository: userCollectionItemRepository,
		CollectionItemRepository:     collectionItemRepository,
	}
}

func (s *collectionService) CollectionList(userID string) (*collectionListResponse, error) {
	if _, err := s.UserRepository.SelectUserByPK(userID); err != nil {
		return nil, err
	}
	userCollectionItems, err := s.UserCollectionItemRepository.SelectUserCollectionItems(userID)
	if err != nil {
		return nil, err
	}
	userCollectionItemIDMap := make(map[string]struct{}, len(userCollectionItems))
	for _, item := range userCollectionItems {
		userCollectionItemIDMap[item.CollectionItemID] = struct{}{}
	}
	masterCollections, err := s.CollectionItemRepository.SelectAllCollectionItems()
	if err != nil {
		return nil, err
	}
	collections := make([]*collection, 0, len(masterCollections))
	for _, masterCollection := range masterCollections {
		_, hasItem := userCollectionItemIDMap[masterCollection.ID]
		collection := &collection{
			CollectionID: masterCollection.ID,
			Name:         masterCollection.Name,
			Rarity:       masterCollection.Rarity,
			HasItem:      hasItem,
		}
		collections = append(collections, collection)
	}
	return &collectionListResponse{
		Collections: collections,
	}, nil
}
