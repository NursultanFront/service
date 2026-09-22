package notificationapp

import (
	"net/http"

	"github.com/ardanlabs/service/business/domain/notificationbus"
	"github.com/ardanlabs/service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	NotificationBus notificationbus.ExtBusiness
}

// Routes adds specific routes for this group.
//
// TODO(dev): this has no auth middleware, unlike auditapp/userapp - decide
// whether this service needs authentication and add mid.Authenticate(...)
// if so.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.NotificationBus)

	app.HandlerFunc(http.MethodGet, version, "/notifications", api.query)
}
