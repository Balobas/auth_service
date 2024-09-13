package deliveryGrpc

import (
	"github.com/balobas/auth_service/pkg/auth_v1"
)

type AuthServerGrpc struct {
	auth_v1.UnimplementedAuthServer

	cfg Config

	ucUsers        UcUsers
	ucAuth         UcAuth
	ucVerification UcVerification
}

type Config interface{

}

func NewAuthServerGRPC(
	cfg Config,
	ucUsers UcUsers,
	ucAuth UcAuth,
	ucVerification UcVerification,
) *AuthServerGrpc {
	return &AuthServerGrpc{
		cfg:            cfg,
		ucUsers:        ucUsers,
		ucAuth:         ucAuth,
		ucVerification: ucVerification,
	}
}
