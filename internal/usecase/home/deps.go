package home

import (
	"context"

	"github.com/ardanlabs/service/internal/models"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/google/uuid"
)

// Storer interface declares the behaviour this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, hme models.Home) error
	Update(ctx context.Context, hme models.Home) error
	Delete(ctx context.Context, hme models.Home) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.Home, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (models.Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]models.Home, error)
}

// userBusiness declares only the user behavior this package calls.
//
// NOTE(AI): narrowed from userbus.ExtBusiness by Claude by direct request
// (layout refactor to Method 1, deps.go on the consumer side).
type userBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (userbus.ExtBusiness, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (models.User, error)
}
