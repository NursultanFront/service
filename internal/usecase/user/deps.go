package user

import (
	"context"
	"net/mail"

	"github.com/ardanlabs/service/internal/models"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/google/uuid"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, usr models.User) error
	Update(ctx context.Context, usr models.User) error
	Delete(ctx context.Context, usr models.User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (models.User, error)
}
