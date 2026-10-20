package models

import "errors"

// Domain errors shared by the usecase and repository layers.
var (
	ErrNotFound = errors.New("home not found")
)
