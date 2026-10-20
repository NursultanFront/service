package main

import (
	"context"
	"os"

	metricsapp "github.com/ardanlabs/service/internal/app/metrics"
)

// NOTE(AI): thin main (migration step (a)); config lives in internal/config/metrics,
// wiring in internal/app/metrics. tag is set at build time via -X main.tag.

var tag = "develop"

func main() {
	if err := metricsapp.Run(context.Background(), tag); err != nil {
		os.Exit(1)
	}
}
