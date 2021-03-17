package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type gameFinishRequest struct {
	Score int `json:"score"`
}
type gameFinishResponse struct {
	Coin int `json:"coin"`
}
type gameHandler struct {
	HttpResponse response.HttpResponseInterface
	GameService  service.GameServiceInterface
}

func NewGameHandler(
	httpResponse response.HttpResponseInterface,
	gameService service.GameServiceInterface,
) *gameHandler {
	return &gameHandler{
		HttpResponse: httpResponse,
		GameService:  gameService,
	}
}

func (h *gameHandler) HandleGameFinish(w http.ResponseWriter, r *http.Request) {
	var finishRequest gameFinishRequest
	if err := json.NewDecoder(r.Body).Decode(&finishRequest); err != nil {
		h.HttpResponse.Failed(w, "failed to decode request body", err, http.StatusBadRequest)
		return
	}
	if finishRequest.Score < 0 {
		h.HttpResponse.Failed(w, fmt.Sprintf("score validation failed. Score=%d", finishRequest.Score), fmt.Errorf("%w. Score=%d", derror.ErrInvalidRequest, finishRequest.Score), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}
	res, err := h.GameService.GameFinish(userID, finishRequest.Score)
	if err != nil {
		h.HttpResponse.Failed(w, "failed to finish game", err, http.StatusInternalServerError)
		return
	}
	transferredResponse := &gameFinishResponse{
		Coin: res.GottenCoin,
	}
	h.HttpResponse.Success(w, transferredResponse)
}
