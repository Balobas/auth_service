package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/balobas/auth_service_bln/pkg/auth_v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServiceGrpc struct {
	auth_v1.UnimplementedAuthServer
}

type Config interface {
}

func NewAuthService(cfg Config) *AuthServiceGrpc {
	return &AuthServiceGrpc{}
}

func (a *AuthServiceGrpc) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	fmt.Printf("%+v\n", req)
	return &auth_v1.CreateResponse{
		Id: 1,
	}, nil
}

func (a *AuthServiceGrpc) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	return &auth_v1.GetResponse{
		Id:        1,
		Name:      "test",
		Email:     "test@email.com",
		Role:      auth_v1.Role_user,
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}, nil
}

func (a *AuthServiceGrpc) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (a *AuthServiceGrpc) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, nil
}
