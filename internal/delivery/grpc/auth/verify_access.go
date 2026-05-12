package deliveryGrpcAuth

import (
	"context"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) VerifyAccess(ctx context.Context, req *auth_v1.VerifyAccessRequest) (*emptypb.Empty, error) {
	log.Printf("authServerGrpc.VerifyAccess: uri %s method %s", req.GetUri(), req.GetMethod())

	if err := s.ucAuth.VerifyAccess(ctx, req.GetUri(), req.GetMethod(), req.GetAccessjwt()); err != nil {
		log.Printf("authServerGrpc.VerifyAccess: failed to verify access: %v", err)
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
