package models

import (
	"context"
	"net/mail"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/google/uuid"
)

// Storer declares the behavior the user usecase needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}
