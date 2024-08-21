package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/balobas/auth_service_bln/internal/config"
	deliveryGrpc "github.com/balobas/auth_service_bln/internal/delivery/grpc"
	repositoryPostgres "github.com/balobas/auth_service_bln/internal/repository/postgres"
	"github.com/balobas/auth_service_bln/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	flag.Parse()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()

	// Инициализация конфигов
	pgConfig := config.NewConfigPG()
	grpcConfig, err := config.NewConfigGRPC()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Инициализация зависимостей
	repo, err := repositoryPostgres.New(ctx, pgConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("successfuly connected to db")

	// Инициализация сервиса
	server := grpc.NewServer()
	reflection.Register(server)
	auth_v1.RegisterAuthServer(server, deliveryGrpc.NewAuthService(nil, repo))

	// Запуск сервиса
	conn, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer conn.Close()

	if err := server.Serve(conn); err != nil {
		log.Printf("failed to serve %v\n", err)
	}
}
