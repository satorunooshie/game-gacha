package middleware

import (
	"context"
	"log"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/model"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx != nil {
			ctx = context.Background()
		}
		token := r.Header.Get("x-token")
		if token == "" {
			log.Println("x-token is empty")
			return
		}
		user, err := model.SelectUserByAuthToken(token)
		if err != nil {
			log.Println(err)
			response.InternalServerError(w, "Invalid token")
			return
		}
		if user == nil {
			log.Printf("user not found. token=%s", token)
			response.BadRequest(w, "Invalid token")
			return
		}
		ctx = dcontext.SetUserID(ctx, user.ID)
		next(w, r.WithContext(ctx))
	}
}
