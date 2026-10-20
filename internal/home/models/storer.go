package models

import (
	"context"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/google/uuid"
)

// Storer declares the behaviour the home usecase needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, hme Home) error
	Update(ctx context.Context, hme Home) error
	Delete(ctx context.Context, hme Home) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Home, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Home, error)
}
