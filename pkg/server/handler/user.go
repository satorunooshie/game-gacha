package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"game-gacha/pkg/dcontext"
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

func HandleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createRequest userCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			log.Println(err, "failed to decode json request")
			response.BadRequest(w, "BadRequest")
			return
		}
		res, err := service.UserCreate(createRequest.Name)
		if err != nil {
			log.Println(err, "failed to create user")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := &userCreateResponse{
			Token: res.Token,
		}
		response.Success(w, transferredResponse)
	}
}
func HandleUserGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println("userID is empty")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		user, err := service.UserGet(userID)
		if err != nil {
			log.Println(err, "failed to get user")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		transferredResponse := &userGetResponse{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
			Coin:      user.Coin,
		}
		response.Success(w, transferredResponse)
	}
}
func HandleUserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateRequest userUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			log.Println(err, "failed to decode json request")
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
		if err := service.UserUpdate(userID, updateRequest.Name); err != nil {
			log.Println(err, "failed to update user")
			response.InternalServerError(w, "Internal Server Error")
			return
		}
		response.Success(w, nil)
	}
}