package main

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server %v", err)
	}
	defer conn.Close()

	client := auth_v1.NewAuthClient(conn)

	resp, err := client.Create(context.Background(), &auth_v1.CreateRequest{
		Name:  "hui",
		Email: "hui@mail.com",
	})
	if err != nil {
		log.Printf("error while create: %v", err)
		return
	}

	log.Printf("%v", resp)
}
