package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type gameFinishRequest struct {
	Score int `json:"score"`
}
type gameFinishResponse struct {
	Coin int `json:"coin"`
}

func HandleGameFinish() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var finishRequest gameFinishRequest
		if err := json.NewDecoder(r.Body).Decode(&finishRequest); err != nil {
			log.Println(err, "failed to decode json request")
			response.BadRequest(w, "Bad Request")
			return
		}
		if finishRequest.Score < 0 {
			log.Println(fmt.Errorf("score validation failed. Score=%d", finishRequest.Score))
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
		res, err := service.GameFinish(userID, finishRequest.Score)
		if err != nil {
			log.Println(err, "failed to finish game")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := &gameFinishResponse{
			Coin: res.GottenCoin,
		}
		response.Success(w, transferredResponse)
	}
}
