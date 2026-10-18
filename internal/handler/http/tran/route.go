package tran

import (
	"net/http"

	"github.com/ardanlabs/service/internal/handler/http/mid"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/sys/auth"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	productbus "github.com/ardanlabs/service/internal/usecase/product"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/ardanlabs/service/pkg/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	DB         *sqlx.DB
	UserBus    userbus.ExtBusiness
	ProductBus productbus.ExtBusiness
	AuthClient authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.UserBus, cfg.ProductBus)

	app.HandlerFunc(http.MethodPost, version, "/tranexample", api.create, authen, ruleAdmin, transaction)
}
