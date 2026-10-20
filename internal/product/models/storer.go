package models

import (
	"context"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/google/uuid"
)

// Storer declares the behavior the product usecase needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}
