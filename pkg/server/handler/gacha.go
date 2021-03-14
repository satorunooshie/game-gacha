package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type gachaDrawRequest struct {
	Times int `json:"times"`
}
type gachaDrawResponse struct {
	Results []*result `json:"results"`
}
type result struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

func HandleGachaDraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var drawRequest gachaDrawRequest
		if err := json.NewDecoder(r.Body).Decode(&drawRequest); err != nil {
			log.Println(err, "failed to decode json request")
			response.BadRequest(w, "Bad Request")
			return
		}
		if drawRequest.Times <= 0 {
			log.Println("request times is invalid")
			response.BadRequest(w, "Bad Request")
			return
		}

		ctx := r.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(w, "Internal Server Error")
			return
		}

		res, err := service.GachaDraw(userID, drawRequest.Times)
		if err != nil {
			if errors.Is(err, derror.ErrCoinShortage) {
				log.Println(err, "failed to draw gacha")
				response.BadRequest(w, "Bad Request")
				return
			}
			log.Println(err, "failed to draw gacha")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := make([]*result, 0, drawRequest.Times)
		for _, v := range res.Results {
			result := result{
				CollectionID: v.CollectionID,
				Name:         v.Name,
				Rarity:       v.Rarity,
				IsNew:        v.IsNew,
			}
			transferredResponse = append(transferredResponse, &result)
		}
		response.Success(w, &gachaDrawResponse{
			Results: transferredResponse,
		})
	}
}
