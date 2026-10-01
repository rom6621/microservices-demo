package domain

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrValidation      = errors.New("validation error")
)
