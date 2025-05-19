package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (s *AuthServerGrpc) CreateAdmin(ctx context.Context, req *auth_v1.AdminCreateRequest) (*auth_v1.AdminCreateResponse, error) {
	log.Printf("authServerGrpc.CreateAdmin")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("authServerGrpc.CreateAdmin: failed to get token context\n")
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "failed to get token from context")
	}

	var token string

	accessJwtSlice := md.Get("accessJwt")
	if len(accessJwtSlice) != 0 {
		token = accessJwtSlice[0]
	}

	uid, err := s.ucUsers.CreateAdmin(
		ctx, entity.User{
			Email: req.GetEmail(),
		},
		req.GetPassword(),
		token,
	)
	if err != nil {
		log.Printf("authServerGrpc.CreateAdmin: failed to create admin user: %v", err)
		return nil, err
	}

	return &auth_v1.AdminCreateResponse{
		Uid: uid.String(),
	}, nil
}
