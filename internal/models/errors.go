package models

import "errors"

// NOTE(AI): moved from internal/usecase/* in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// Domain errors shared by the usecase and repository layers.
var (
	ErrHomeNotFound    = errors.New("home not found")
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserUniqueEmail = errors.New("email is not unique")
)
