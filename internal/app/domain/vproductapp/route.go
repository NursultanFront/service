package vproductapp

import (
	"net/http"

	"github.com/ardanlabs/service/internal/app/sdk/auth"
	"github.com/ardanlabs/service/internal/app/sdk/authclient"
	"github.com/ardanlabs/service/internal/app/sdk/mid"
	"github.com/ardanlabs/service/internal/business/domain/userbus"
	"github.com/ardanlabs/service/internal/business/domain/vproductbus"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log         *logger.Logger
	UserBus     userbus.ExtBusiness
	VProductBus vproductbus.ExtBusiness
	AuthClient  authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.VProductBus)

	app.HandlerFunc(http.MethodGet, version, "/vproducts", api.query, authen, ruleAdmin)
}
