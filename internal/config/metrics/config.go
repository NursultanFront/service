// Package metrics holds the configuration of the metrics service.
package metrics

import (
	"time"

	"github.com/ardanlabs/conf/v3"
)

// NOTE(AI): moved from cmd/metrics/main.go in migration step (a); same fields and defaults.

// Config holds the METRICS_* settings of the metrics service.
type Config struct {
	conf.Version
	Web struct {
		DebugHost string `conf:"default:0.0.0.0:4010"`
	}
	Expvar struct {
		Host            string        `conf:"default:0.0.0.0:4000"`
		Route           string        `conf:"default:/metrics"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Prometheus struct {
		Host            string        `conf:"default:0.0.0.0:4020"`
		Route           string        `conf:"default:/metrics"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Collect struct {
		From string `conf:"default:http://localhost:3010/debug/vars"`
	}
	Publish struct {
		To       string        `conf:"default:console"`
		Interval time.Duration `conf:"default:5s"`
	}
}

const prefix = "METRICS"

// Parse reads the configuration from command-line flags and METRICS_*
// environment variables. When help is requested it returns conf.ErrHelpWanted
// together with the help text.
func Parse(build string) (Config, string, error) {
	cfg := Config{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		return Config{}, help, err
	}

	return cfg, "", nil
}
