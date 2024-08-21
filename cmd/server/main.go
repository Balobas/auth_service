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

	fmt.Println(configPath)
	if err := config.Load(configPath); err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()

	grpcConfig, err := config.NewConfigGRPC()
	if err != nil {
		log.Fatalf("%v", err)
	}

	conn, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer conn.Close()

	repo, err := repositoryPostgres.New(ctx, config.NewConfigPG())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("successfuly connected to db")

	server := grpc.NewServer()
	reflection.Register(server)

	auth_v1.RegisterAuthServer(server, deliveryGrpc.NewAuthService(nil, repo))
	log.Printf("server listening at %s\n", conn.Addr())

	if err := server.Serve(conn); err != nil {
		log.Printf("failed to serve %v\n", err)
	}
}
