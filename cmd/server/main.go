package main

import (
	"log"
	"net"

	deliveryGrpc "github.com/balobas/auth_service_bln/internal/delivery/grpc"
	"github.com/balobas/auth_service_bln/pkg/auth_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = "localhost:50051"
)

func main() {
	conn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	defer conn.Close()

	server := grpc.NewServer()
	reflection.Register(server)

	auth_v1.RegisterAuthServer(server, deliveryGrpc.NewAuthService(nil))
	log.Printf("server listening at %s\n", conn.Addr())

	if err := server.Serve(conn); err != nil {
		log.Printf("failed to serve %v\n", err)
	}
}
