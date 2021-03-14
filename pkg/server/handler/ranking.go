package handler

import (
	"log"
	"net/http"
	"strconv"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type rankingListResponse struct {
	Ranks []*rank `json:"ranks"`
}
type rank struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

func HandleRankingList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		if start == "" {
			log.Println("Empty Request")
			response.BadRequest(w, "Bad Request")
			return
		}
		startPosition, err := strconv.Atoi(start)
		if err != nil || startPosition <= 0 {
			log.Println(err, "start position is invalid")
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
		res, err := service.RankingList(userID, startPosition, constant.RankingLimit)
		if err != nil {
			log.Println(err, "failed to get ranking list")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := make([]*rank, 0, len(res.Ranks))
		for _, v := range res.Ranks {
			rank := &rank{
				UserID:   v.UserID,
				UserName: v.UserName,
				Rank:     v.Rank,
				Score:    v.Score,
			}
			transferredResponse = append(transferredResponse, rank)
		}
		response.Success(w, &rankingListResponse{
			Ranks: transferredResponse,
		})
	}
}
