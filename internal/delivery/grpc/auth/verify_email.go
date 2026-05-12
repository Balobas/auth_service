package deliveryGrpcAuth

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) VerifyEmail(ctx context.Context, req *auth_v1.VerifyEmailRequest) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.VerifyEmail: token %s", req.GetToken())
	return nil, s.ucVerification.Verify(ctx, req.GetToken())
}
