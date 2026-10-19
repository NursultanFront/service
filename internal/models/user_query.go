package models

import (
	"net/mail"
	"time"

	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/user/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// UserQueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type UserQueryFilter struct {
	ID               *uuid.UUID
	Name             *name.Name
	Email            *mail.Address
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}

// Set of fields that user results can be ordered by.
const (
	UserOrderByID      = "a"
	UserOrderByName    = "b"
	UserOrderByEmail   = "c"
	UserOrderByRoles   = "d"
	UserOrderByEnabled = "e"
)
