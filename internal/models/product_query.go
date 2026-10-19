package models

import (
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/product/{filter,order}.go in migration step (b2) so that
// internal/repository no longer imports internal/usecase; the usecase package
// keeps aliases, so its API is unchanged.

// ProductQueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type ProductQueryFilter struct {
	ID       *uuid.UUID
	Name     *name.Name
	Cost     *float64
	Quantity *int
}

// Set of fields that product results can be ordered by.
const (
	ProductOrderByProductID = "a"
	ProductOrderByUserID    = "b"
	ProductOrderByName      = "c"
	ProductOrderByCost      = "d"
	ProductOrderByQuantity  = "e"
)
