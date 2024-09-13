package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.JwtResponse, error) {
	accessjwt, refreshJwt, err := s.ucAuth.Login(
		ctx, entity.LoginParams{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  accessjwt,
		RefreshJwt: refreshJwt,
	}, nil
}
