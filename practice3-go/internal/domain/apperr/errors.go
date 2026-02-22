package apperr

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
)