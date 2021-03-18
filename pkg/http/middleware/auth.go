package middleware

import (
	"context"
	"fmt"
	"net/http"

	"game-gacha/pkg/dcontext"
	"game-gacha/pkg/derror"
	"game-gacha/pkg/http/response"
	"game-gacha/pkg/server/model"
)

type middleware struct {
	HttpResponse   response.HttpResponseInterface
	UserRepository model.UserRepositoryInterface
}
type MiddlewareInterface interface {
	Authenticate(next http.HandlerFunc) http.HandlerFunc
}

var _ MiddlewareInterface = (*middleware)(nil)

func NewMiddleware(
	httpResponse response.HttpResponseInterface,
	userRepository model.UserRepositoryInterface,
) *middleware {
	return &middleware{
		HttpResponse:   httpResponse,
		UserRepository: userRepository,
	}
}

func (m *middleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx != nil {
			ctx = context.Background()
		}
		token := r.Header.Get("x-token")
		if token == "" {
			m.HttpResponse.Failed(w, "x-token is empty", derror.ErrEmptyToken, http.StatusBadRequest)
			return
		}
		user, err := m.UserRepository.SelectUserByAuthToken(token)
		if err != nil {
			m.HttpResponse.Failed(w, "Invalid token", err, http.StatusBadRequest)
			return
		}
		if user == nil {
			m.HttpResponse.Failed(w, fmt.Sprintf("Invalid token. user not found. token=%s", token), err, http.StatusBadRequest)
			return
		}
		ctx = dcontext.SetUserID(ctx, user.ID)
		next(w, r.WithContext(ctx))
	}
}
