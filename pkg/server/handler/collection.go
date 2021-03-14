package handler

import (
	"fmt"
	"log"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type collectionListRequest struct{}
type collectionListResponse struct {
	Collections []*collection `json:"collections"`
}
type collection struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

func HandleCollectionList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		res, err := service.CollectionList(userID)
		if err != nil {
			log.Println(err, fmt.Sprintf("failed to get user's collection list %s", userID))
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := make([]*collection, 0, len(res.Collections))
		for _, v := range res.Collections {
			collection := &collection{
				CollectionID: v.CollectionID,
				Name:         v.Name,
				Rarity:       v.Rarity,
				HasItem:      v.HasItem,
			}
			transferredResponse = append(transferredResponse, collection)
		}
		response.Success(w, &collectionListResponse{
			Collections: transferredResponse,
		})
	}
}
