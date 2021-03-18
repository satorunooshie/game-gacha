package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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
type gachaHandler struct {
	HttpResponse response.HttpResponseInterface
	GachaService service.GachaServiceInterface
}
type GachaHandlerInterface interface {
	HandleGachaDraw(w http.ResponseWriter, r *http.Request)
}

var _ GachaHandlerInterface = (*gachaHandler)(nil)

func NewGachaHandler(
	httpResponse response.HttpResponseInterface,
	gachaService service.GachaServiceInterface,
) GachaHandlerInterface {
	return &gachaHandler{
		HttpResponse: httpResponse,
		GachaService: gachaService,
	}
}

func (h *gachaHandler) HandleGachaDraw(w http.ResponseWriter, r *http.Request) {
	var drawRequest gachaDrawRequest
	if err := json.NewDecoder(r.Body).Decode(&drawRequest); err != nil {
		h.HttpResponse.Failed(w, "failed to decode request body", err, http.StatusBadRequest)
		return
	}
	if drawRequest.Times <= 0 {
		h.HttpResponse.Failed(w, fmt.Sprintf("request times is invalid. request times=%d", drawRequest.Times), fmt.Errorf("%w. request times=%d", derror.ErrInvalidRequest, drawRequest.Times), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}

	res, err := h.GachaService.GachaDraw(userID, drawRequest.Times)
	if err != nil {
		if errors.Is(err, derror.ErrCoinShortage) {
			h.HttpResponse.Failed(w, "failed to draw gacha", err, http.StatusBadRequest)
			return
		}
		h.HttpResponse.Failed(w, "failed to draw gacha", err, http.StatusInternalServerError)
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
	h.HttpResponse.Success(w, &gachaDrawResponse{
		Results: transferredResponse,
	})
}
