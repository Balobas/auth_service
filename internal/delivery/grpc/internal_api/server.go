package deliveryGrpcAuthInternalApi

import (
	"github.com/balobas/auth_service/pkg/auth_internal_api"
)

type AuthInternalApiServerGrpc struct {
	auth_internal_api.UnimplementedAuthInternalApiServer

	cfg Config

	ucAuth     UcAuth
	ucServices UcServices
}

func NewAuthInternalApiServerGrpc(
	cfg Config,
	ucAuth UcAuth,
	ucServices UcServices,
) *AuthInternalApiServerGrpc {
	return &AuthInternalApiServerGrpc{
		cfg:        cfg,
		ucAuth:     ucAuth,
		ucServices: ucServices,
	}
}
