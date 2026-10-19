// Package home provides business access to home domain.
package home

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/ardanlabs/service/internal/models"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/delegate"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = models.ErrHomeNotFound // NOTE(AI): moved to models in step (b2)
	ErrUserDisabled = errors.New("user disabled")
)

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, nh models.NewHome) (models.Home, error)
	Update(ctx context.Context, hme models.Home, uh models.UpdateHome) (models.Home, error)
	Delete(ctx context.Context, hme models.Home) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.Home, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (models.Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]models.Home, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for home api access.
type Business struct {
	log        *logger.Logger
	userBus    userBusiness
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs a home business API for use.
func NewBusiness(log *logger.Logger, userBus userBusiness, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		userBus:    userBus,
		delegate:   delegate,
		storer:     storer,
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

// NewWithTx constructs a new domain value that will use the
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

// Create adds a new home to the system.
func (b *Business) Create(ctx context.Context, nh models.NewHome) (models.Home, error) {
	usr, err := b.userBus.QueryByID(ctx, nh.UserID)
	if err != nil {
		return models.Home{}, fmt.Errorf("user.querybyid: %s: %w", nh.UserID, err)
	}

	if !usr.Enabled {
		return models.Home{}, ErrUserDisabled
	}

	now := time.Now()

	hme := models.Home{
		ID:   uuid.New(),
		Type: nh.Type,
		Address: models.Address{
			Address1: nh.Address.Address1,
			Address2: nh.Address.Address2,
			ZipCode:  nh.Address.ZipCode,
			City:     nh.Address.City,
			State:    nh.Address.State,
			Country:  nh.Address.Country,
		},
		UserID:      nh.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, hme); err != nil {
		return models.Home{}, fmt.Errorf("create: %w", err)
	}

	return hme, nil
}

// Update modifies information about a home.
func (b *Business) Update(ctx context.Context, hme models.Home, uh models.UpdateHome) (models.Home, error) {
	if uh.Type != nil {
		hme.Type = *uh.Type
	}

	if uh.Address != nil {
		if uh.Address.Address1 != nil {
			hme.Address.Address1 = *uh.Address.Address1
		}

		if uh.Address.Address2 != nil {
			hme.Address.Address2 = *uh.Address.Address2
		}

		if uh.Address.ZipCode != nil {
			hme.Address.ZipCode = *uh.Address.ZipCode
		}

		if uh.Address.City != nil {
			hme.Address.City = *uh.Address.City
		}

		if uh.Address.State != nil {
			hme.Address.State = *uh.Address.State
		}

		if uh.Address.Country != nil {
			hme.Address.Country = *uh.Address.Country
		}
	}

	hme.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, hme); err != nil {
		return models.Home{}, fmt.Errorf("update: %w", err)
	}

	return hme, nil
}

// Delete removes the specified home.
func (b *Business) Delete(ctx context.Context, hme models.Home) error {
	if err := b.storer.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing homes.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]models.Home, error) {
	hmes, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}

// Count returns the total number of homes.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the home by the specified ID.
func (b *Business) QueryByID(ctx context.Context, homeID uuid.UUID) (models.Home, error) {
	hme, err := b.storer.QueryByID(ctx, homeID)
	if err != nil {
		return models.Home{}, fmt.Errorf("query: homeID[%s]: %w", homeID, err)
	}

	return hme, nil
}

// QueryByUserID finds the homes by a specified User ID.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]models.Home, error) {
	hmes, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}
