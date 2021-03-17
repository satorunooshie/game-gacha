package derror

import (
	"errors"
	"fmt"
)

var (
	ErrCoinShortage = errors.New("coin is not enough")
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
