package deliveryGrpcAuthInternalApi

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_internal_api"
)

func (s *AuthInternalApiServerGrpc) ServiceRefresh(ctx context.Context, req *auth_internal_api.ServiceRefreshRequest) (*auth_internal_api.JwtResponse, error) {
	log.Printf("AuthInternalApiServerGrpc.ServiceRefresh")
	accessJwt, refreshJwt, err := s.ucAuth.ServiceRefresh(ctx, req.GetRefreshJwt())
	if err != nil {
		log.Printf("AuthInternalApiServerGrpc.ServiceRefresh: failed: %v", err)
		return nil, err
	}

	return &auth_internal_api.JwtResponse{
		AccessJwt:  accessJwt,
		RefreshJwt: refreshJwt,
	}, nil
}
