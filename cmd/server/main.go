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

	if err := app.NewApp(configPath).Run(ctx); err != nil {
		log.Fatal(err)
	}
}
