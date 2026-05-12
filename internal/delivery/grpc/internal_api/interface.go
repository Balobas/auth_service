package deliveryGrpcAuthInternalApi

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
)

type Config interface {
}

type UcAuth interface {
	ServiceLogin(ctx context.Context, params entity.ServiceLoginParams) (string, string, error)
	ServiceRefresh(ctx context.Context, token string) (string, string, error)
}

type UcServices interface {
	Register(ctx context.Context, req entity.RegisterServiceRequest) error
}
