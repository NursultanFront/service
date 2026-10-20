package audit

import (
	"context"

	auditmodels "github.com/ardanlabs/service/internal/audit/models"

	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, audit auditmodels.Audit) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]auditmodels.Audit, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}
