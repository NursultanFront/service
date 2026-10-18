// Package sales holds the configuration of the sales service.
package sales

import (
	"time"

	"github.com/ardanlabs/conf/v3"
)

// NOTE(AI): moved from cmd/sales/main.go in migration step (a); same fields and defaults.

// Config holds the SALES_* settings of the sales service.
type Config struct {
	conf.Version
	Web struct {
		ReadTimeout        time.Duration `conf:"default:5s"`
		WriteTimeout       time.Duration `conf:"default:10s"`
		IdleTimeout        time.Duration `conf:"default:120s"`
		ShutdownTimeout    time.Duration `conf:"default:20s"`
		APIHost            string        `conf:"default:0.0.0.0:3000"`
		DebugHost          string        `conf:"default:0.0.0.0:3010"`
		CORSAllowedOrigins []string      `conf:"default:*"`
	}
	Auth struct {
		Host      string `conf:"default:http://auth-service:6000"`
		Mechanism string `conf:"default:http"`
		GRPC      struct {
			Host string `conf:"default:auth-service:6001"`
		}
	}
	DB struct {
		User         string `conf:"default:postgres"`
		Password     string `conf:"default:postgres,mask"`
		Host         string `conf:"default:database-service"`
		Name         string `conf:"default:postgres"`
		MaxIdleConns int    `conf:"default:0"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
	Redis struct {
		Host     string `conf:"default:redis:6379"`
		Password string `conf:"mask"`
		DB       int    `conf:"default:0"`
	}
	Tempo struct {
		Host        string  `conf:"default:tempo:4317"`
		ServiceName string  `conf:"default:sales"`
		Probability float64 `conf:"default:0.05"`
		// Shouldn't use a high Probability value in non-developer systems.
		// 0.05 should be enough for most systems. Some might want to have
		// this even lower.
	}
}

const prefix = "SALES"

// Parse reads the configuration from command-line flags and SALES_*
// environment variables. When help is requested it returns conf.ErrHelpWanted
// together with the help text.
func Parse(build string) (Config, string, error) {
	cfg := Config{
		Version: conf.Version{
			Build: build,
			Desc:  "Sales",
		},
	}

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		return Config{}, help, err
	}

	return cfg, "", nil
}
