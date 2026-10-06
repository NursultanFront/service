//go:build reporting

package build

import (
	"github.com/ardanlabs/service/internal/app/domain/checkapp"
	"github.com/ardanlabs/service/internal/app/domain/vproductapp"
	"github.com/ardanlabs/service/internal/app/sdk/mux"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

// Routes binds the reporting routes for the sales service.
func Routes() rpt {
	return rpt{}
}

type rpt struct{}

// Add implements the RouterAdder interface.
func (rpt) Add(app *web.App, cfg mux.Config) {
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	vproductapp.Routes(app, vproductapp.Config{
		UserBus:     cfg.BusConfig.UserBus,
		VProductBus: cfg.BusConfig.VProductBus,
		AuthClient:  cfg.SalesConfig.AuthClient,
	})
}
