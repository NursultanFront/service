// Package usecase provides business access to view product domain.
package usecase

import (
	"context"
	"fmt"
	"slices"

	vproductmodels "github.com/ardanlabs/service/internal/vproduct/models"

	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
)

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]vproductmodels.Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for view product access.
type Business struct {
	storer Storer
}

// NewBusiness constructs a vproduct business API for use.
func NewBusiness(storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		storer: storer,
	})

	for _, ext := range slices.Backward(extensions) {

		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// Query retrieves a list of existing products.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]vproductmodels.Product, error) {
	users, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// Count returns the total number of products.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}
