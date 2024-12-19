package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) CreateAdmin(ctx context.Context, req *auth_v1.AdminCreateRequest) (*auth_v1.AdminCreateResponse, error) {
	log.Printf("auth.CreateAdmin\n")

	uid, err := s.ucUsers.CreateAdmin(
		ctx, entity.User{
			Email: req.GetEmail(),
		},
		req.GetPassword(),
		req.GetJwt(),
	)
	if err != nil {
		log.Printf("failed to create admin user\n")
		return nil, err
	}

	return &auth_v1.AdminCreateResponse{
		Uid: uid.String(),
	}, nil
}
