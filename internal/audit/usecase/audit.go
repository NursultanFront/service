// Package usecase provides a business logic layer for handling audit events.
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	auditmodels "github.com/ardanlabs/service/internal/audit/models"

	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/google/uuid"
)

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	Create(ctx context.Context, na auditmodels.NewAudit) (auditmodels.Audit, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]auditmodels.Audit, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for audit access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a audit business API for use.
func NewBusiness(log *logger.Logger, storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		log:    log,
		storer: storer,
	})

	for _, ext := range slices.Backward(extensions) {

		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// Create adds a new audit record to the system.
func (b *Business) Create(ctx context.Context, na auditmodels.NewAudit) (auditmodels.Audit, error) {
	jsonData, err := json.Marshal(na.Data)
	if err != nil {
		return auditmodels.Audit{}, fmt.Errorf("marshal object: %w", err)
	}

	audit := auditmodels.Audit{
		ID:        uuid.New(),
		ObjID:     na.ObjID,
		ObjDomain: na.ObjDomain,
		ObjName:   na.ObjName,
		ActorID:   na.ActorID,
		Action:    na.Action,
		Data:      jsonData,
		Message:   na.Message,
		Timestamp: time.Now(),
	}

	if err := b.storer.Create(ctx, audit); err != nil {
		return auditmodels.Audit{}, fmt.Errorf("create audit: %w", err)
	}

	return audit, nil
}

// Query retrieves a list of existing audit records.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]auditmodels.Audit, error) {
	audits, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query audits: %w", err)
	}

	return audits, nil
}

// Count returns the total number of users.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}
