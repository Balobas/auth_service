package deliveryGrpcAuth

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) Verify(ctx context.Context, req *auth_v1.VerifyRequest) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.Verify: ")
	_, err := s.ucAuth.VerifyAuth(ctx, req.GetJwt())
	return nil, err
}
