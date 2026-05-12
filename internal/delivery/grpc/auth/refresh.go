package deliveryGrpcAuth

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
)

func (s *AuthServerGrpc) Refresh(ctx context.Context, req *auth_v1.RefreshRequest) (*auth_v1.JwtResponse, error) {
	log.Printf("authServerGrpc.Refresh")
	accessJwt, refreshJwt, err := s.ucAuth.Refresh(ctx, req.GetRefreshJwt())
	if err != nil {
		log.Printf("authServerGrpc.Refresh: failed: %v", err)
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  accessJwt,
		RefreshJwt: refreshJwt,
	}, nil
}
