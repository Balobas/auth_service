package grpc

import (
	"context"
	"fmt"

	converterGrpc "github.com/balobas/auth_service/internal/delivery/grpc/converter"
	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServerGrpc struct {
	auth_v1.UnimplementedAuthServer

	usersService UsersService
}

type UsersService interface {
	Create(ctx context.Context, user entity.User) (int64, error)
	Get(ctx context.Context, id int64) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id int64) error
}

type Config interface{}

func NewAuthServerGRPC(cfg Config, usersService UsersService) *AuthServerGrpc {
	return &AuthServerGrpc{
		usersService: usersService,
	}
}

func (a *AuthServerGrpc) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	fmt.Printf("create request: %+v\n", req)

	id, err := a.usersService.Create(ctx, converterGrpc.FromCreateUserRequestToUserEntity(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}

func (a *AuthServerGrpc) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	fmt.Printf("get user request: %+v\n", req)

	user, err := a.usersService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converterGrpc.FromUserEntityToGetResponse(user), nil
}

func (a *AuthServerGrpc) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	return nil, a.usersService.Update(ctx, converterGrpc.FromUpdateRequestToUserEntity(req))
}

func (a *AuthServerGrpc) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	return nil, a.usersService.Delete(ctx, req.GetId())
}
