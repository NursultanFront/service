// Package useraudit provides an extension for userbus that adds
// auditing functionality.
package useraudit

import (
	"context"
	"net/mail"

	"github.com/ardanlabs/service/internal/models"
	"github.com/ardanlabs/service/internal/models/types/domain"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	auditbus "github.com/ardanlabs/service/internal/usecase/audit"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/google/uuid"
)

// Extension provides a wrapper for audit functionality around the userbus.
type Extension struct {
	bus      userbus.ExtBusiness
	auditBus auditbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the userbus with audit.
func NewExtension(auditBus auditbus.ExtBusiness) userbus.Extension {
	return func(bus userbus.ExtBusiness) userbus.ExtBusiness {
		return &Extension{
			bus:      bus,
			auditBus: auditBus,
		}
	}
}

// NewWithTx does not apply auditing.
func (ext *Extension) NewWithTx(tx sqldb.CommitRollbacker) (userbus.ExtBusiness, error) {
	return ext.bus.NewWithTx(tx)
}

// Create applies auditing to the user creation process.
func (ext *Extension) Create(ctx context.Context, actorID uuid.UUID, nu models.NewUser) (models.User, error) {
	usr, err := ext.bus.Create(ctx, actorID, nu)
	if err != nil {
		return models.User{}, err
	}

	na := models.NewAudit{
		ObjID:     usr.ID,
		ObjDomain: domain.User,
		ObjName:   usr.Name,
		ActorID:   actorID,
		Action:    "created",
		Data:      nu,
		Message:   "user created",
	}

	if _, err := ext.auditBus.Create(ctx, na); err != nil {
		return models.User{}, err
	}

	return usr, nil
}

// Update applies auditing to the user update process.
func (ext *Extension) Update(ctx context.Context, actorID uuid.UUID, usr models.User, uu models.UpdateUser) (models.User, error) {
	usr, err := ext.bus.Update(ctx, actorID, usr, uu)
	if err != nil {
		return models.User{}, err
	}

	na := models.NewAudit{
		ObjID:     usr.ID,
		ObjDomain: domain.User,
		ObjName:   usr.Name,
		ActorID:   actorID,
		Action:    "updated",
		Data:      uu,
		Message:   "user updated",
	}

	if _, err := ext.auditBus.Create(ctx, na); err != nil {
		return models.User{}, err
	}

	return usr, nil
}

// Delete applies auditing to the user deletion process.
func (ext *Extension) Delete(ctx context.Context, actorID uuid.UUID, usr models.User) error {
	if err := ext.bus.Delete(ctx, actorID, usr); err != nil {
		return err
	}

	na := models.NewAudit{
		ObjID:     usr.ID,
		ObjDomain: domain.User,
		ObjName:   usr.Name,
		ActorID:   actorID,
		Action:    "deleted",
		Data:      nil,
		Message:   "user deleted",
	}

	if _, err := ext.auditBus.Create(ctx, na); err != nil {
		return err
	}

	return nil
}

// Query does not apply auditing.
func (ext *Extension) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, page page.Page) ([]models.User, error) {
	return ext.bus.Query(ctx, filter, orderBy, page)
}

// Count does not apply auditing.
func (ext *Extension) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	return ext.bus.Count(ctx, filter)
}

// QueryByID does not apply auditing.
func (ext *Extension) QueryByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	return ext.bus.QueryByID(ctx, userID)
}

// QueryByEmail does not apply auditing.
func (ext *Extension) QueryByEmail(ctx context.Context, email mail.Address) (models.User, error) {
	return ext.bus.QueryByEmail(ctx, email)
}

// Authenticate does not apply auditing.
func (ext *Extension) Authenticate(ctx context.Context, email mail.Address, password string) (models.User, error) {
	return ext.bus.Authenticate(ctx, email, password)
}
