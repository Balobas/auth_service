package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	deliveryGrpc "github.com/balobas/auth_service_bln/internal/delivery/grpc"
	repositoryPostgres "github.com/balobas/auth_service_bln/internal/repository/postgres"
	"github.com/balobas/auth_service_bln/pkg/auth_v1"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = "[::]:50051"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

type Config struct {
	dbName     string
	pgUser     string
	pgPassword string
	pgPort     string
	pgHost     string
}

func NewConfig() Config {
	return Config{
		dbName:     os.Getenv("PG_DATABASE_NAME"),
		pgUser:     os.Getenv("PG_USER"),
		pgPassword: os.Getenv("PG_PASSWORD"),
		pgPort:     os.Getenv("PG_PORT"),
		pgHost:     os.Getenv("PG_HOST"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.pgHost, c.pgPort, c.dbName, c.pgUser, c.pgPassword,
	)
}

func main() {
	flag.Parse()

	fmt.Println(configPath)
	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("%v", err)
	}

	ctx := context.Background()

	conn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer conn.Close()

	config := NewConfig()

	repo, err := repositoryPostgres.New(ctx, config)
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
