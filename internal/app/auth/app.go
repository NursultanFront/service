// Package auth wires every dependency of the auth service and runs it.
package auth

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	authcfg "github.com/ardanlabs/service/internal/config/auth"
	grpcauthapp "github.com/ardanlabs/service/internal/handler/grpc/auth"
	mux "github.com/ardanlabs/service/internal/handler/http/server"
	userdb "github.com/ardanlabs/service/internal/repository/user"
	usercache "github.com/ardanlabs/service/internal/repository/user/cache"
	"github.com/ardanlabs/service/internal/sys/auth"
	sqldb "github.com/ardanlabs/service/internal/sys/db"
	"github.com/ardanlabs/service/internal/sys/debug"
	"github.com/ardanlabs/service/internal/sys/delegate"
	"github.com/ardanlabs/service/internal/sys/rediscache"
	userbus "github.com/ardanlabs/service/internal/usecase/user"
	"github.com/ardanlabs/service/pkg/foundation/keystore"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/ardanlabs/service/pkg/foundation/otel"
)

// NOTE(AI): moved from cmd/auth/main.go in migration step (a); wiring only,
// no logic changes. cmd/auth/main.go now just calls Run.

// Run builds the logger and every dependency of the auth service and blocks
// until the service shuts down. Startup errors are logged before being returned.
func Run(ctx context.Context, tag string) error {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", otel.GetTraceID, events)

	if err := run(ctx, log, tag); err != nil {
		log.Error(ctx, "startup", "err", err)
		return err
	}

	return nil
}

func run(ctx context.Context, log *logger.Logger, tag string) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg, help, err := authcfg.Parse(tag)
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
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), redisClient, time.Minute))

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	// Check the environment first to see if a key is being provided. Then
	// load any private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.

	ks := keystore.New()

	n1, err := ks.LoadByJSON(cfg.Auth.KeysJSON)
	if err != nil {
		return fmt.Errorf("loading keys by env: %w", err)
	}

	n2, err := ks.LoadByFileSystem(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("loading keys by fs: %w", err)
	}

	if n1+n2 == 0 {
		return errors.New("no keys exist")
	}

	authCfg := auth.Config{
		Log:       log,
		UserBus:   userBus,
		KeyLookup: ks,
		Issuer:    cfg.Auth.Issuer,
	}

	ath := auth.New(authCfg)

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
			UserBus: userBus,
		},
		AuthConfig: mux.AuthConfig{
			Auth: ath,
		},
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfgMux, Routes(), mux.WithCORS(cfg.Web.CORSAllowedOrigins)),
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
	// Start gRPC Service

	log.Info(ctx, "startup", "status", "initializing gRPC support")

	lis, err := net.Listen("tcp", cfg.Web.GRPCHost)
	if err != nil {
		return fmt.Errorf("failed to listen on host %s : %w", cfg.Web.GRPCHost, err)
	}

	grpcService := grpcauthapp.Start(ctx, grpcauthapp.Config{
		Log:      log,
		Auth:     ath,
		Listener: lis,
		UserBus:  userBus,
	})

	defer grpcService.Shutdown()

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
