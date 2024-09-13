package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {
	log.Printf("auth.Register\n")

	uid, err := s.ucUsers.Register(
		ctx, entity.User{
			Email: req.GetEmail(),
		},
		req.GetPassword(),
	)
	if err != nil {
		log.Printf("failed to register user\n")
		return nil, err
	}

	return &auth_v1.RegisterResponse{
		Uid: uid.String(),
	}, nil
}
