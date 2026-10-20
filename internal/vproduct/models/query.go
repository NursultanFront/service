package models

import (
	"github.com/ardanlabs/service/internal/types/name"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/vproduct/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID       *uuid.UUID
	Name     *name.Name
	Cost     *float64
	Quantity *int
	UserName *name.Name
}

// Set of fields that vproduct results can be ordered by.
const (
	OrderByProductID = "a"
	OrderByUserID    = "b"
	OrderByName      = "c"
	OrderByCost      = "d"
	OrderByQuantity  = "e"
	OrderByUserName  = "f"
)
