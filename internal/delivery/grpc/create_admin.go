package deliveryGrpc

import (
	"context"
	"errors"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/grpc/metadata"
)

func (s *AuthServerGrpc) CreateAdmin(ctx context.Context, req *auth_v1.AdminCreateRequest) (*auth_v1.AdminCreateResponse, error) {
	log.Printf("auth.CreateAdmin\n")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("failed to get token context\n")
		return nil, errors.New("failed to get token context")
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
		log.Printf("failed to create admin user\n")
		return nil, err
	}

	return &auth_v1.AdminCreateResponse{
		Uid: uid.String(),
	}, nil
}
