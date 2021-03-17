package derror

import (
	"errors"
	"fmt"
)

var (
	ErrCoinShortage   = errors.New("coin is not enough")
	ErrUserNotFound   = errors.New("user not found")
	ErrEmptyUserID    = errors.New("userID is empty")
	ErrEmptyToken     = errors.New("x-token is empty")
	ErrEmptyRequest   = errors.New("request body is empty")
	ErrInvalidRequest = errors.New("request is invalid")
)

type ApplicationError struct {
	Message string
	Err     error
	Code    int
}

func (ae ApplicationError) Error() string {
	if ae.Err != nil {
		return fmt.Sprintf("%s: %s", ae.Message, ae.Err.Error())
	}
	return ae.Message
}
