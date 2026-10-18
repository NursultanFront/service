package main

import (
	"context"
	"embed"
	"os"

	salesapp "github.com/ardanlabs/service/internal/app/sales"
)

// NOTE(AI): thin main (migration step (a)); config lives in internal/config/sales,
// wiring in internal/app/sales. tag is set at build time via -X main.tag.

//go:embed static
var static embed.FS

var tag = "develop"

func main() {
	if err := salesapp.Run(context.Background(), tag, static); err != nil {
		os.Exit(1)
	}
}
