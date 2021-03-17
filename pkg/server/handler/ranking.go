package handler

import (
	"net/http"
	"strconv"

	"game-gacha/pkg/constant"
	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
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
type rankingHandler struct {
	HttpResponse   response.HttpResponseInterface
	RankingService service.RankingServiceInterface
}

func NewRankingHandler(
	httpResponse response.HttpResponseInterface,
	rankingService service.RankingServiceInterface,
) *rankingHandler {
	return &rankingHandler{
		HttpResponse:   httpResponse,
		RankingService: rankingService,
	}
}

func (h *rankingHandler) HandleRankingList(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	if start == "" {
		h.HttpResponse.Failed(w, "empty request", derror.ErrEmptyRequest, http.StatusBadRequest)
		return
	}
	startPosition, err := strconv.Atoi(start)
	if err != nil || startPosition <= 0 {
		h.HttpResponse.Failed(w, "start position is invalid", err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}
	res, err := h.RankingService.RankingList(userID, startPosition, constant.RankingLimit)
	if err != nil {
		h.HttpResponse.Failed(w, "failed to get ranking list", err, http.StatusInternalServerError)
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
	h.HttpResponse.Success(w, &rankingListResponse{
		Ranks: transferredResponse,
	})
}
