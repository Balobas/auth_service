package main

import (
	"context"
	"flag"
	"log"
  
	"github.com/balobas/auth_service/internal/app"
	"github.com/balobas/auth_service/internal/client/pg"
	"github.com/balobas/auth_service/internal/config"
	deliveryGrpc "github.com/balobas/auth_service/internal/delivery/grpc"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
	usersService "github.com/balobas/auth_service/internal/service/users"
	"github.com/balobas/auth_service/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
