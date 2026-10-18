package vproduct

import (
	"context"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.VProduct, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}
