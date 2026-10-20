package models

import "errors"

// Domain errors shared by the usecase and repository layers.
var (
	ErrNotFound    = errors.New("user not found")
	ErrUniqueEmail = errors.New("email is not unique")
)
