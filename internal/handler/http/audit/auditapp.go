// Package audit maintains the app layer api for the audit domain.
package audit

import (
	"context"
	"net/http"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/handler/http/query"
	"github.com/ardanlabs/service/internal/sys/order"
	"github.com/ardanlabs/service/internal/sys/page"
	auditbus "github.com/ardanlabs/service/internal/usecase/audit"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

type app struct {
	auditBus auditbus.ExtBusiness
}

func newApp(auditBus auditbus.ExtBusiness) *app {
	return &app{
		auditBus: auditBus,
	}
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, auditbus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	adts, err := a.auditBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	total, err := a.auditBus.Count(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppAudits(adts), total, page)
}
