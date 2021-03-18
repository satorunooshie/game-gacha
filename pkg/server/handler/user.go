package handler

import (
	"encoding/json"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/service"
)

type userCreateRequest struct {
	Name string `json:"name"`
}
type userCreateResponse struct {
	Token string `json:"token"`
}
type userGetRequest struct{}
type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
	Coin      int    `json:"coin"`
}
type userUpdateRequest struct {
	Name string `json:"name"`
}
type userUpdateResponse struct{}

type userHandler struct {
	HttpResponse response.HttpResponseInterface
	UserService  service.UserServiceInterface
}
type UserHandlerInterface interface {
	HandleUserCreate(w http.ResponseWriter, r *http.Request)
	HandleUserGet(w http.ResponseWriter, r *http.Request)
	HandleUserUpdate(w http.ResponseWriter, r *http.Request)
}

var _ UserHandlerInterface = (*userHandler)(nil)

func NewUserHandler(
	httpResponse response.HttpResponseInterface,
	userService service.UserServiceInterface,
) UserHandlerInterface {
	return &userHandler{
		HttpResponse: httpResponse,
		UserService: userService,
	}
}
func (h *userHandler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	var createRequest userCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		h.HttpResponse.Failed(w, "failed to decode request body", err, http.StatusBadRequest)
		return
	}
	res, err := h.UserService.UserCreate(createRequest.Name)
	if err != nil {
		h.HttpResponse.Failed(w, "failed to create user", err, http.StatusInternalServerError)
		return
	}
	transferredResponse := &userCreateResponse{
		Token: res.Token,
	}
	h.HttpResponse.Success(w, transferredResponse)
}
func (h *userHandler) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}
	user, err := h.UserService.UserGet(userID)
	if err != nil {
		h.HttpResponse.Failed(w, "failed to get user", err, http.StatusInternalServerError)
		return
	}
	transferredResponse := &userGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	}
	h.HttpResponse.Success(w, transferredResponse)
}
func (h *userHandler) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	var updateRequest userUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		h.HttpResponse.Failed(w, "failed to decode request body", err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userID := dcontext.GetUserIDFromContext(ctx)
	if userID == "" {
		h.HttpResponse.Failed(w, "userID is empty", derror.ErrEmptyUserID, http.StatusInternalServerError)
		return
	}
	if err := h.UserService.UserUpdate(userID, updateRequest.Name); err != nil {
		h.HttpResponse.Failed(w, "failed to update user", err, http.StatusInternalServerError)
		return
	}
	h.HttpResponse.Success(w, nil)
}
