package handler

import (
	"fmt"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
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

type collectionHandler struct {
	HttpResponse      response.HttpResponseInterface
	CollectionService service.CollectionServiceInterface
}
type CollectionHandlerInterface interface {
	HandleCollectionList(w http.ResponseWriter, r *http.Request)
}

var _ CollectionHandlerInterface = (*collectionHandler)(nil)

func NewCollectionHandler(
	httpResponse response.HttpResponseInterface,
	collectionService service.CollectionServiceInterface,
) CollectionHandlerInterface {
	return &collectionHandler{
		HttpResponse:      httpResponse,
		CollectionService: collectionService,
	}
}
func (h *collectionHandler) HandleCollectionList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}
	res, err := h.CollectionService.CollectionList(userID)
	if err != nil {
		h.HttpResponse.Failed(w, fmt.Sprintf("failed to get user's collection list %s", userID), err, http.StatusInternalServerError)
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
	h.HttpResponse.Success(w, &collectionListResponse{
		Collections: transferredResponse,
	})
}
