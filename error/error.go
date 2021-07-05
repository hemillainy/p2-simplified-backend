package error

import "errors"

var (
	ErrInvalidTransaction = errors.New("insufficient balance")
)
