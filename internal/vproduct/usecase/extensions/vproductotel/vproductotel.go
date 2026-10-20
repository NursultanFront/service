// Package vproductotel provides an extension for vproductbus that adds
// otel tracking.
package vproductotel

import (
	"context"

	vproductmodels "github.com/ardanlabs/service/internal/vproduct/models"

	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	vproductbus "github.com/ardanlabs/service/internal/vproduct/usecase"
	"github.com/ardanlabs/service/pkg/foundation/otel"
)

// Extension provides a wrapper for otel functionality around the vproductbus.
type Extension struct {
	bus vproductbus.ExtBusiness
}

// NewExtension constructs a new extension that wraps the vproductbus with otel.
func NewExtension() vproductbus.Extension {
	return func(bus vproductbus.ExtBusiness) vproductbus.ExtBusiness {
		return &Extension{
			bus: bus,
		}
	}
}

func (ext *Extension) Query(ctx context.Context, filter vproductbus.QueryFilter, orderBy order.By, page page.Page) ([]vproductmodels.Product, error) {
	ctx, span := otel.AddSpan(ctx, "business.vproductbus.query")
	defer span.End()

	usr, err := ext.bus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (ext *Extension) Count(ctx context.Context, filter vproductbus.QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.vproductbus.count")
	defer span.End()

	usr, err := ext.bus.Count(ctx, filter)
	if err != nil {
		return 0, err
	}

	return usr, nil
}
