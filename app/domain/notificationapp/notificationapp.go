// Package notificationapp maintains the app layer api for the notification
// domain.
package notificationapp

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/business/domain/notificationbus"
	"github.com/ardanlabs/service/foundation/web"
)

type app struct {
	notificationBus notificationbus.ExtBusiness
}

func newApp(notificationBus notificationbus.ExtBusiness) *app {
	return &app{
		notificationBus: notificationBus,
	}
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	ns, err := a.notificationBus.Query(ctx)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	return toAppNotifications(ns)
}
