package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.JwtResponse, error) {
	log.Printf("authServerGrpc.Login: email %s", req.GetEmail())
	accessjwt, refreshJwt, err := s.ucAuth.Login(
		ctx, entity.LoginParams{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		log.Printf("authServerGrpc.Login: failed to login user with email %s: %v", req.GetEmail(), err)
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  accessjwt,
		RefreshJwt: refreshJwt,
	}, nil
}
