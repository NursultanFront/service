package models

import (
	"time"

	"github.com/ardanlabs/service/internal/models/types/money"
	"github.com/ardanlabs/service/internal/models/types/name"
	"github.com/ardanlabs/service/internal/models/types/quantity"
	"github.com/google/uuid"
)

// VProduct represents an individual product with extended information, including the owning user's name.
type VProduct struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
	UserName    name.Name
}
