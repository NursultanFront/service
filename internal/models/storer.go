package models

import (
	"context"
	"net/mail"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/google/uuid"
)

// NOTE(AI): moved from internal/usecase/{user,product,home}/deps.go in migration
// step (b2) so that internal/repository no longer imports internal/usecase; the
// usecase packages keep aliases, so their API is unchanged.

// HomeStorer declares the behaviour the home usecase needs to persist and
// retrieve data.
type HomeStorer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (HomeStorer, error)
	Create(ctx context.Context, hme Home) error
	Update(ctx context.Context, hme Home) error
	Delete(ctx context.Context, hme Home) error
	Query(ctx context.Context, filter HomeQueryFilter, orderBy order.By, page page.Page) ([]Home, error)
	Count(ctx context.Context, filter HomeQueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Home, error)
}

// ProductStorer declares the behavior the product usecase needs to persist and
// retrieve data.
type ProductStorer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ProductStorer, error)
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter ProductQueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter ProductQueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// UserStorer declares the behavior the user usecase needs to persist and
// retrieve data.
type UserStorer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (UserStorer, error)
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter UserQueryFilter, orderBy order.By, page page.Page) ([]User, error)
	Count(ctx context.Context, filter UserQueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}
