package derror

import (
	"errors"
)

var (
	ErrCoinShortage = errors.New("coin is not enough")
)
