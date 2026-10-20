//go:build !crud && !reporting

package sales

import (
	auditapp "github.com/ardanlabs/service/internal/handler/http/audit"
	checkapp "github.com/ardanlabs/service/internal/handler/http/check"
	homeapp "github.com/ardanlabs/service/internal/handler/http/home"
	productapp "github.com/ardanlabs/service/internal/handler/http/product"
	rawapp "github.com/ardanlabs/service/internal/handler/http/raw"
	tranapp "github.com/ardanlabs/service/internal/handler/http/tran"
	userapp "github.com/ardanlabs/service/internal/handler/http/user"
	vproductapp "github.com/ardanlabs/service/internal/handler/http/vproduct"
	mux "github.com/ardanlabs/service/internal/server/http"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

// Routes binds all the routes for the sales service.
func Routes() all {
	return all{}
}

type all struct{}

// Add implements the RouterAdder interface.
func (all) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homeapp.Routes(app, homeapp.Config{
		Log:        cfg.Log,
		HomeBus:    cfg.BusConfig.HomeBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	productapp.Routes(app, productapp.Config{
		Log:        cfg.Log,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	rawapp.Routes(app)

	tranapp.Routes(app, tranapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		ProductBus: cfg.BusConfig.ProductBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	auditapp.Routes(app, auditapp.Config{
		Log:        cfg.Log,
		AuditBus:   cfg.BusConfig.AuditBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})

	vproductapp.Routes(app, vproductapp.Config{
		Log:         cfg.Log,
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.SalesConfig.AuthClient,
	})
}
