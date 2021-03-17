package service

import (
	"database/sql"
	"fmt"
	"math/rand"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/db"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/server/model"
)

type gachaDrawResponse struct {
	Results []*result
}
type result struct {
	CollectionID string
	Name         string
	Rarity       int
	IsNew        bool
}
type collectionItem struct {
	Name   string
	Rarity int
}
type gachaService struct {
	UserRepository               model.UserRepositoryInterface
	GachaProbabilityRepository   model.GachaProbabilityRepositoryInterface
	UserCollectionItemRepository model.UserCollectionItemRepositoryInterface
	CollectionItemRepository     model.CollectionItemRepositoryInterface
}
type GachaServiceInterface interface {
	GachaDraw(userID string, times int) (*gachaDrawResponse, error)
}

func NewGachaService(
	userRepository model.UserRepositoryInterface,
	gachaProbabilityRepository model.GachaProbabilityRepositoryInterface,
	userCollectionItemRepository model.UserCollectionItemRepositoryInterface,
	collectionItemRepository model.CollectionItemRepositoryInterface,
) *gachaService {
	return &gachaService{
		UserRepository:               userRepository,
		GachaProbabilityRepository:   gachaProbabilityRepository,
		UserCollectionItemRepository: userCollectionItemRepository,
		CollectionItemRepository:     collectionItemRepository,
	}
}

func (s *gachaService) GachaDraw(userID string, times int) (*gachaDrawResponse, error) {
	results := make([]*result, 0, times)
	if err := db.DB.Transaction(func(tx *sql.Tx) error {
		user, err := s.UserRepository.SelectUserByPKForUpdate(tx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("%w. userID=%s", derror.ErrUserNotFound, userID)
		}
		remainingCoins := user.Coin - (times * constant.GachaCoinConsumption)
		if remainingCoins < 0 {
			return fmt.Errorf("%w. shortage=%d", derror.ErrCoinShortage, remainingCoins)
		}
		userCollectionItems, err := s.UserCollectionItemRepository.SelectUserCollectionItems(userID)
		if err != nil {
			return err
		}
		userCollectionItemIDMap := make(map[string]struct{}, len(userCollectionItems))
		for _, item := range userCollectionItems {
			userCollectionItemIDMap[item.CollectionItemID] = struct{}{}
		}
		collections, err := s.CollectionItemRepository.SelectAllCollectionItems()
		if err != nil {
			return err
		}
		collectionMap := make(map[string]*collectionItem, len(collections))
		for _, collection := range collections {
			collectionMap[collection.ID] = &collectionItem{
				Name:   collection.Name,
				Rarity: collection.Rarity,
			}
		}
		gachaProbabilities, err := s.GachaProbabilityRepository.SelectGachaProbabilities()
		if err != nil {
			return err
		}
		sum := 0
		for _, gachaProbability := range gachaProbabilities {
			sum += gachaProbability.Ratio
		}
		newItemIDs := make([]string, 0, times)
		for i := times; i > 0; i-- {
			random := rand.Intn(sum)
			var (
				tmpRate        int
				selectedItemID string
				isNew          bool
			)
			for _, gachaProbability := range gachaProbabilities {
				tmpRate += gachaProbability.Ratio
				if tmpRate > random {
					selectedItemID = gachaProbability.CollectionItemID
					break
				}
			}
			if _, ok := userCollectionItemIDMap[selectedItemID]; !ok {
				isNew = true
				newItemIDs = append(newItemIDs, selectedItemID)
				userCollectionItemIDMap[selectedItemID] = struct{}{}
			}
			result := result{
				CollectionID: selectedItemID,
				Name:         collectionMap[selectedItemID].Name,
				Rarity:       collectionMap[selectedItemID].Rarity,
				IsNew:        isNew,
			}
			results = append(results, &result)
		}
		user.Coin = remainingCoins
		if len(newItemIDs) > 0 {
			if err = s.UserCollectionItemRepository.SaveUserCollectionItems(tx, newItemIDs, userID); err != nil {
				return err
			}
		}
		if err = s.UserRepository.UpdateUserCoinByPK(tx, user.Coin, userID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &gachaDrawResponse{
		Results: results,
	}, nil
}
