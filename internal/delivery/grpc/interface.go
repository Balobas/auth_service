package deliveryGrpc

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type UcUsers interface {
	Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error)
	CreateAdmin(ctx context.Context, user entity.User, password string, token string) (uuid.UUID, error)
	GetAdmins(ctx context.Context, token string) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, password string) error
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	DeleteUser(ctx context.Context, userUid uuid.UUID) error
}

type UcAuth interface {
	Login(ctx context.Context, params entity.LoginParams) (string, string, error)
	Logout(ctx context.Context, userUid uuid.UUID) error
	Refresh(ctx context.Context, token string) (string, string, error)
	VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error)
	UpdateUserCreds(ctx context.Context, user entity.User, password string) (string, string, error)
}

type UcVerification interface {
	Verify(ctx context.Context, token string) error
}

type UcPermissions interface {
	CreatePermission(ctx context.Context, perm entity.Permission) error
	DeletePermission(ctx context.Context, key string) error
	GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error)
	UpdatePermission(ctx context.Context, perm entity.Permission) error
	AddPermissionToUser(ctx context.Context, userUid uuid.UUID, permKey string) error
	RemoveUserPermission(ctx context.Context, userUid uuid.UUID, permKey string) error
}
