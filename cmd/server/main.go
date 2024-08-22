package main

import (
	"context"
	"flag"
	"log"

	"github.com/balobas/auth_service_bln/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	app, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
