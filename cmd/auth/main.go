package main

import (
	"context"
	"os"

	authapp "github.com/ardanlabs/service/internal/app/auth"
)

// NOTE(AI): thin main (migration step (a)); config lives in internal/config/auth,
// wiring in internal/app/auth. tag is set at build time via -X main.tag.

var tag = "develop"

func main() {
	if err := authapp.Run(context.Background(), tag); err != nil {
		os.Exit(1)
	}
}
