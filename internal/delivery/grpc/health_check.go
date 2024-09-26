package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthServerGrpc) HealthCheck(context.Context, *emptypb.Empty) (*auth_v1.HealthCheckResponse, error) {
	return &auth_v1.HealthCheckResponse{
		Status: "OK",
	}, nil
}
