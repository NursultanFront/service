package audit

import (
	"net/http"

	"github.com/ardanlabs/service/internal/handler/http/mid"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/sys/auth"
	auditbus "github.com/ardanlabs/service/internal/usecase/audit"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	AuditBus   auditbus.ExtBusiness
	AuthClient authclient.Authenticator
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.AuditBus)

	app.HandlerFunc(http.MethodGet, version, "/audits", api.query, authen, ruleAdmin)
}
