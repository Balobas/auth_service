package deliveryGrpcInterceptors

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
)

type UcAuth interface {
	VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error)
}
