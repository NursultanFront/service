// Package usecase provides business access to product domain.
package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	productmodels "github.com/ardanlabs/service/internal/product/models"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/delegate"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = productmodels.ErrNotFound // NOTE(AI): moved to models in step (b2)
	ErrUserDisabled = errors.New("user disabled")
	ErrInvalidCost  = errors.New("cost not valid")
)

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, np productmodels.NewProduct) (productmodels.Product, error)
	Update(ctx context.Context, prd productmodels.Product, up productmodels.UpdateProduct) (productmodels.Product, error)
	Delete(ctx context.Context, prd productmodels.Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]productmodels.Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (productmodels.Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]productmodels.Product, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for product access.
type Business struct {
	log        *logger.Logger
	userBus    userBusiness
	storer     Storer
	delegate   *delegate.Delegate
	extensions []Extension
}

// NewBusiness constructs a product business API for use.
func NewBusiness(log *logger.Logger, userBus userBusiness, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		userBus:    userBus,
		storer:     storer,
		delegate:   delegate,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for _, ext := range slices.Backward(extensions) {

		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, userBus, b.delegate, storer, b.extensions...)

	return nb, nil
}

// Create adds a new product to the system.
func (b *Business) Create(ctx context.Context, np productmodels.NewProduct) (productmodels.Product, error) {
	usr, err := b.userBus.QueryByID(ctx, np.UserID)
	if err != nil {
		return productmodels.Product{}, fmt.Errorf("user.querybyid: %s: %w", np.UserID, err)
	}

	if !usr.Enabled {
		return productmodels.Product{}, ErrUserDisabled
	}

	now := time.Now()

	prd := productmodels.Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, prd); err != nil {
		return productmodels.Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// Update modifies information about a product.
func (b *Business) Update(ctx context.Context, prd productmodels.Product, up productmodels.UpdateProduct) (productmodels.Product, error) {
	if up.Name != nil {
		prd.Name = *up.Name
	}

	if up.Cost != nil {
		prd.Cost = *up.Cost
	}

	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}

	prd.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, prd); err != nil {
		return productmodels.Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

// Delete removes the specified product.
func (b *Business) Delete(ctx context.Context, prd productmodels.Product) error {
	if err := b.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing products.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]productmodels.Product, error) {
	prds, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count returns the total number of products.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the product by the specified ID.
func (b *Business) QueryByID(ctx context.Context, productID uuid.UUID) (productmodels.Product, error) {
	prd, err := b.storer.QueryByID(ctx, productID)
	if err != nil {
		return productmodels.Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID finds the products by a specified User ID.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]productmodels.Product, error) {
	prds, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
