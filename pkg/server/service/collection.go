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

func CollectionList(userID string) (*collectionListResponse, error) {
	if _, err := model.SelectUserByPK(userID); err != nil {
		return nil, err
	}
	userCollectionItems, err := model.SelectUserCollectionItems(userID)
	if err != nil {
		return nil, err
	}
	userCollectionItemIDMap := make(map[string]struct{}, len(userCollectionItems))
	for _, item := range userCollectionItems {
		userCollectionItemIDMap[item.CollectionItemID] = struct{}{}
	}
	masterCollections, err := model.SelectAllCollectionItems()
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
