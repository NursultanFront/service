// This is the starting point for the notifications service. It exposes a
// small REST API over notifications that are (eventually) populated by
// consuming Kafka events published by other services (see productkafka).
//
// TODO(dev): this is a skeleton - see the TODOs below and in the packages
// it wires together (notificationbus, notificationmem, kafkaclient).
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/service/app/domain/notificationapp"
	"github.com/ardanlabs/service/app/sdk/debug"
	"github.com/ardanlabs/service/app/sdk/mid"
	"github.com/ardanlabs/service/business/domain/notificationbus"
	"github.com/ardanlabs/service/business/domain/notificationbus/stores/notificationmem"
	"github.com/ardanlabs/service/business/sdk/kafkaclient"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/otel"
	"github.com/ardanlabs/service/foundation/web"
)

var tag = "develop"

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "NOTIFICATIONS", traceIDFn)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:5000"`
			DebugHost       string        `conf:"default:0.0.0.0:5010"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
		Kafka struct {
			Brokers []string `conf:"default:kafka:19092"`
			Topic   string   `conf:"default:products"`
			GroupID string   `conf:"default:notifications"`
		}
		Tempo struct {
			Host        string  `conf:"default:tempo:4317"`
			ServiceName string  `conf:"default:notifications"`
			Probability float64 `conf:"default:0.05"`
		}
	}{
		Build: tag,
		Desc:  "copyright information here",
	}

	const prefix = "NOTIFICATIONS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", tag)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "err", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start Tracing Support

	traceProvider, teardown, err := otel.InitTracing(otel.Config{
		ServiceName: cfg.Tempo.ServiceName,
		Host:        cfg.Tempo.Host,
		Probability: cfg.Tempo.Probability,
	})
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer teardown(ctx)

	tracer := traceProvider.Tracer(cfg.Tempo.ServiceName)

	// -------------------------------------------------------------------------
	// Create Business Packages
	//
	// TODO(dev): notificationmem.NewStore() is a no-op placeholder. Swap it
	// for a real store once you've decided how notifications are persisted.

	notificationBus := notificationbus.NewBusiness(log, notificationmem.NewStore())

	// -------------------------------------------------------------------------
	// Start Kafka Consumer
	//
	// TODO(dev): this reader is constructed but nothing reads from it yet.
	// Launch a goroutine (see foundation/worker, or a plain `go func`) that
	// loops on reader.ReadMessage(ctx), unmarshals the event, and calls
	// notificationBus.Create. Remember to reader.Close() on shutdown.

	reader := kafkaclient.NewReader(kafkaclient.Config{Brokers: cfg.Kafka.Brokers}, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	defer reader.Close()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	app := web.NewApp(
		log.Info,
		tracer,
		mid.Otel(tracer),
		mid.Logger(log),
		mid.Errors(log),
		mid.Metrics(),
		mid.Panics(),
	)

	notificationapp.Routes(app, notificationapp.Config{
		NotificationBus: notificationBus,
	})

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      app,
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

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
