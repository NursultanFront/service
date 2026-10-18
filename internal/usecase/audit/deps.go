package audit

import (
	"context"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, audit models.Audit) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.Audit, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}
