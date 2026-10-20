// Package models holds the vproduct domain entities, query filter, and order-by
// fields shared by its usecase and repository layers.
package models

import (
	"time"

	"github.com/ardanlabs/service/internal/types/money"
	"github.com/ardanlabs/service/internal/types/name"
	"github.com/ardanlabs/service/internal/types/quantity"
	"github.com/google/uuid"
)

// Product represents an individual product with extended information, including the owning user's name.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
	UserName    name.Name
}
