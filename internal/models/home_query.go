package models

import (
	"time"

	"github.com/ardanlabs/service/internal/models/types/home"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/home/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// HomeQueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type HomeQueryFilter struct {
	ID               *uuid.UUID
	UserID           *uuid.UUID
	Type             *home.Home
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}

// Set of fields that home results can be ordered by.
const (
	HomeOrderByID     = "a"
	HomeOrderByType   = "b"
	HomeOrderByUserID = "c"
)
