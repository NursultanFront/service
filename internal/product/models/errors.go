package models

import "errors"

// Domain errors shared by the usecase and repository layers.
var (
	ErrNotFound = errors.New("product not found")
)
