package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Refresh(ctx context.Context, req *auth_v1.RefreshRequest) (*auth_v1.JwtResponse, error) {
	accessJwt, refreshJwt, err := s.ucAuth.Refresh(ctx, req.GetRefreshJwt())
	if err != nil {
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  accessJwt,
		RefreshJwt: refreshJwt,
	}, nil
}
