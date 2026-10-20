// Package userotel provides an extension for userbus that adds
// otel tracking.
package userotel

import (
	"context"
	"net/mail"

	usermodels "github.com/ardanlabs/service/internal/user/models"

	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/ardanlabs/service/pkg/foundation/otel"
	"github.com/google/uuid"
)

// Extension provides a wrapper for otel functionality around the userbus.
type Extension struct {
	bus userbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the userbus with otel.
func NewExtension() userbus.Extension {
	return func(bus userbus.ExtBusiness) userbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

// NewWithTx does not apply otel.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (userbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Create applies otel to the user creation process.
func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, nu usermodels.NewUser) (usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.create")
	defer span.End()

	usr, err := ext.bus.Create(ctx, actorID, nu)
	if err != nil {
		return usermodels.User{}, err
	}

	return usr, nil
}

// Update applies otel to the user update process.
func (ext *Extension) Update(ctx context.Context, actorID uuid.UUID, usr usermodels.User, uu usermodels.UpdateUser) (usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.update")
	defer span.End()

	usr, err := ext.bus.Update(ctx, actorID, usr, uu)
	if err != nil {
		return usermodels.User{}, err
	}

	return usr, nil
}

// Delete applies otel to the user deletion process.
func (ext *Extension) Delete(ctx context.Context, actorID uuid.UUID, usr usermodels.User) error {
	ctx, span := otel.AddSpan(ctx, "business.userbus.delete")
	defer span.End()

	if err := ext.bus.Delete(ctx, actorID, usr); err != nil {
		return err
	}

	return nil
}

// Query applies otel to the user query process.
func (ext *Extension) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, page page.Page) ([]usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.query")
	defer span.End()

	return ext.bus.Query(ctx, filter, orderBy, page)
}

// Count applies otel to the user count process.
func (ext *Extension) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.count")
	defer span.End()

	return ext.bus.Count(ctx, filter)
}

// QueryByID applies otel to the user query by ID process.
func (ext *Extension) QueryByID(ctx context.Context, userID uuid.UUID) (usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.querybyid")
	defer span.End()

	return ext.bus.QueryByID(ctx, userID)
}

// QueryByEmail applies otel to the user query by email process.
func (ext *Extension) QueryByEmail(ctx context.Context, email mail.Address) (usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.querybyemail")
	defer span.End()

	return ext.bus.QueryByEmail(ctx, email)
}

// Authenticate applies otel to the user authentication process.
func (ext *Extension) Authenticate(ctx context.Context, email mail.Address, password string) (usermodels.User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.authenticate")
	defer span.End()

	return ext.bus.Authenticate(ctx, email, password)
}
