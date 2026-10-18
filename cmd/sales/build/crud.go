//go:build crud

package build

import (
	auditapp "github.com/ardanlabs/service/internal/handler/http/audit"
	checkapp "github.com/ardanlabs/service/internal/handler/http/check"
	homeapp "github.com/ardanlabs/service/internal/handler/http/home"
	productapp "github.com/ardanlabs/service/internal/handler/http/product"
	mux "github.com/ardanlabs/service/internal/handler/http/server"
	tranapp "github.com/ardanlabs/service/internal/handler/http/tran"
	userapp "github.com/ardanlabs/service/internal/handler/http/user"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

// Routes binds the crud routes for the sales service.
func Routes() crud {
	return crud{}
}

type crud struct{}

// Add implements the RouterAdder interface.
func (crud) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	tranapp.Routes(app, tranapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		Log:        cfg.Log,
		AuthClient: cfg.SalesConfig.AuthClient,
		DB:         cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})
}
