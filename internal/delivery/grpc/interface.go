package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type UcUsers interface {
	Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user entity.User, password string) error
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, bool, error)
	DeleteUser(ctx context.Context, userUid uuid.UUID) error
}

type UcAuth interface {
	Login(ctx context.Context, params entity.LoginParams) (string, string, error)
	Logout(ctx context.Context, userUid uuid.UUID) error
	Refresh(ctx context.Context, token string) (string, string, error)
	VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error)
}

type UcVerification interface {
	Verify(ctx context.Context, token string) error
}
