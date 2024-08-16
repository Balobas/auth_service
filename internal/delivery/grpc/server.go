package grpc

import (
	"context"
	"fmt"
	"time"

	repositoryPostgres "github.com/balobas/auth_service_bln/internal/repository/postgres"
	"github.com/balobas/auth_service_bln/pkg/auth_v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServiceGrpc struct {
	auth_v1.UnimplementedAuthServer

	repository *repositoryPostgres.Repository
}

type Config interface{}

func NewAuthService(cfg Config, repo *repositoryPostgres.Repository) *AuthServiceGrpc {
	return &AuthServiceGrpc{
		repository: repo,
	}
}

func (a *AuthServiceGrpc) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	fmt.Printf("create request: %+v\n", req)

	id, err := a.repository.CreateUser(ctx, req, req.Role)
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}

func (a *AuthServiceGrpc) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	fmt.Printf("get user request: %+v\n", req)

	user, err := a.repository.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      auth_v1.Role(auth_v1.Role_value[user.Role]),
		CreatedAt: timestamppb.New(time.UnixMicro(user.CreatedAt)),
		UpdatedAt: timestamppb.New(time.UnixMicro(user.UpdatedAt)),
	}, nil
}

func (a *AuthServiceGrpc) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, a.repository.UpdateUser(ctx, req.GetId(), req.GetName().GetValue(), req.GetEmail().GetValue())
}

func (a *AuthServiceGrpc) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, a.repository.DeleteUser(ctx, req.GetId())
}
