// Package sales wires every dependency of the sales service and runs it.
package sales

import (
	"context"
	"embed"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	salescfg "github.com/ardanlabs/service/internal/config/sales"
	mux "github.com/ardanlabs/service/internal/handler/http/server"
	auditdb "github.com/ardanlabs/service/internal/repository/audit"
	homedb "github.com/ardanlabs/service/internal/repository/home"
	productpg "github.com/ardanlabs/service/internal/repository/product/postgres"
	userdb "github.com/ardanlabs/service/internal/repository/user"
	usercache "github.com/ardanlabs/service/internal/repository/user/cache"
	vproductdb "github.com/ardanlabs/service/internal/repository/vproduct"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/service/authclient/grpc"
	http2 "github.com/ardanlabs/service/internal/service/authclient/http"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/debug"
	"github.com/ardanlabs/service/internal/sys/delegate"
	"github.com/ardanlabs/service/internal/sys/rediscache"
	auditbus "github.com/ardanlabs/service/internal/usecase/audit"
	"github.com/ardanlabs/service/internal/usecase/audit/extensions/auditotel"
	homebus "github.com/ardanlabs/service/internal/usecase/home"
	"github.com/ardanlabs/service/internal/usecase/home/extensions/homeotel"
	productbus "github.com/ardanlabs/service/internal/usecase/product"
	"github.com/ardanlabs/service/internal/usecase/product/extensions/productotel"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/ardanlabs/service/internal/usecase/user/extensions/useraudit"
	"github.com/ardanlabs/service/internal/usecase/user/extensions/userotel"
	vproductbus "github.com/ardanlabs/service/internal/usecase/vproduct"
	"github.com/ardanlabs/service/internal/usecase/vproduct/extensions/vproductotel"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/ardanlabs/service/pkg/foundation/otel"
)

// NOTE(AI): moved from cmd/sales/main.go in migration step (a); wiring only,
// no logic changes. cmd/sales/main.go now just calls Run.

/*
	Need to figure out timeouts for http service.
*/

// Run builds the logger and every dependency of the sales service and blocks
// until the service shuts down. Startup errors are logged before being returned.
func Run(ctx context.Context, tag string, static embed.FS) error {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", otel.GetTraceID, events)

	if err := run(ctx, log, tag, static); err != nil {
		log.Error(ctx, "startup", "err", err)
		return err
	}

	return nil
}

func run(ctx context.Context, log *logger.Logger, tag string, static embed.FS) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg, help, err := salescfg.Parse(tag)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer db.Close()

	// -------------------------------------------------------------------------
	// Redis Support

	log.Info(ctx, "startup", "status", "initializing redis support", "host", cfg.Redis.Host)

	redisClient, err := rediscache.Open(rediscache.Config{
		Host:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("connecting to redis: %w", err)
	}

	defer redisClient.Close()

	// -------------------------------------------------------------------------
	// Create Business Packages

	delegate := delegate.New(log)

	auditOtelExt := auditotel.NewExtension()
	auditStorage := auditdb.NewStore(log, db)
	auditBus := auditbus.NewBusiness(log, auditStorage, auditOtelExt)

	userOtelExt := userotel.NewExtension()
	userAuditExt := useraudit.NewExtension(auditBus)
	userStorage := usercache.NewStore(log, userdb.NewStore(log, db), redisClient, time.Minute)
	userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

	productOtelExt := productotel.NewExtension()
	productStorage := productpg.NewStore(log, db)
	productBus := productbus.NewBusiness(log, userBus, delegate, productStorage, productOtelExt)

	homeOtelExt := homeotel.NewExtension()
	homeStorage := homedb.NewStore(log, db)
	homeBus := homebus.NewBusiness(log, userBus, delegate, homeStorage, homeOtelExt)

	vproductOtelExt := vproductotel.NewExtension()
	vproductStorage := vproductdb.NewStore(log, db)
	vproductBus := vproductbus.NewBusiness(vproductStorage, vproductOtelExt)

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	var authClient authclient.Authenticator
	switch cfg.Auth.Mechanism {
	case "grpc":
		authClient, err = grpc.New(log, cfg.Auth.GRPC.Host)
	default:
		authClient, err = http2.New(log, cfg.Auth.Host)
	}

	if err != nil {
		log.Error(ctx, "failed to initialize authentication client", "error", err)
		return fmt.Errorf("failed to initialize authentication client: %w", err)
	}

	defer authClient.Close()

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing tracing support")

	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		ServiceName: cfg.Tempo.ServiceName,
		Host:        cfg.Tempo.Host,
		ExcludedRoutes: map[string]struct{}{
			"/v1/liveness":  {},
			"/v1/readiness": {},
		},
		Probability: cfg.Tempo.Probability,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}

	defer teardown(context.Background())

	tracer := traceProvider.Tracer(cfg.Tempo.ServiceName)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Build:  tag,
		Log:    log,
		DB:     db,
		Tracer: tracer,
		BusConfig: mux.BusConfig{
			AuditBus:    auditBus,
			UserBus:     userBus,
			ProductBus:  productBus,
			HomeBus:     homeBus,
			VProductBus: vproductBus,
		},
		SalesConfig: mux.SalesConfig{
			AuthClient: authClient,
		},
	}

	webAPI := mux.WebAPI(cfgMux,
		Routes(),
		mux.WithCORS(cfg.Web.CORSAllowedOrigins),
		mux.WithFileServer(false, static, "static", "/"),
	)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      webAPI,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
