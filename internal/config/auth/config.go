// Package auth holds the configuration of the auth service.
package auth

import (
	"time"

	"github.com/ardanlabs/conf/v3"
)

// NOTE(AI): moved from cmd/auth/main.go in migration step (a); same fields and defaults.

// Config holds the AUTH_* settings of the auth service.
type Config struct {
	conf.Version
	Web struct {
		ReadTimeout        time.Duration `conf:"default:5s"`
		WriteTimeout       time.Duration `conf:"default:10s"`
		IdleTimeout        time.Duration `conf:"default:120s"`
		ShutdownTimeout    time.Duration `conf:"default:20s"`
		APIHost            string        `conf:"default:0.0.0.0:6000"`
		DebugHost          string        `conf:"default:0.0.0.0:6010"`
		GRPCHost           string        `conf:"default:0.0.0.0:6001"`
		CORSAllowedOrigins []string      `conf:"default:*"`
	}
	Auth struct {
		KeysJSON   string `conf:"mask"`
		KeysFolder string `conf:"default:zarf/keys/"`
		ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
		Issuer     string `conf:"default:service project"`
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
		ServiceName string  `conf:"default:auth"`
		Probability float64 `conf:"default:0.05"`
	}
}

const prefix = "AUTH"

// Parse reads the configuration from command-line flags and AUTH_*
// environment variables. When help is requested it returns conf.ErrHelpWanted
// together with the help text.
func Parse(build string) (Config, string, error) {
	cfg := Config{
		Version: conf.Version{
			Build: build,
			Desc:  "Auth",
		},
	}

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		return Config{}, help, err
	}

	return cfg, "", nil
}
