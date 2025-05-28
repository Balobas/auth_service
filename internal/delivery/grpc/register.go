package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.RegisterResponse, error) {
	log.Printf("authServerGrpc.Register: email %s", req.GetEmail())

	uid, err := s.ucUsers.Register(
		ctx, entity.User{
			Email: req.GetEmail(),
		},
		req.GetPassword(),
	)
	if err != nil {
		log.Printf("authServerGrpc.Register: failed to register user with email %s: %v", req.GetEmail(), err)
		return nil, err
	}

	return &auth_v1.RegisterResponse{
		Uid: uid.String(),
	}, nil
}
