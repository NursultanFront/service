package models

import (
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/vproduct/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// VProductQueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type VProductQueryFilter struct {
	ID       *uuid.UUID
	Name     *name.Name
	Cost     *float64
	Quantity *int
	UserName *name.Name
}

// Set of fields that vproduct results can be ordered by.
const (
	VProductOrderByProductID = "a"
	VProductOrderByUserID    = "b"
	VProductOrderByName      = "c"
	VProductOrderByCost      = "d"
	VProductOrderByQuantity  = "e"
	VProductOrderByUserName  = "f"
)
