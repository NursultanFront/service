package vproduct

import (
	"context"

	vproductmodels "github.com/ardanlabs/service/internal/vproduct/models"

	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]vproductmodels.Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}
