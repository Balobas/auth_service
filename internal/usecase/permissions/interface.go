package useCasePermissions

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type PermissionsRepository interface {
	CreatePermission(ctx context.Context, perm entity.Permission) error
	UpdatePermission(ctx context.Context, perm entity.Permission) error
	GetPermission(ctx context.Context, key string) (entity.Permission, bool, error)
	GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error)
	DeletePermission(ctx context.Context, key string) error
	RemovePermissionFromUsers(ctx context.Context, key string, limitUsers int64) ([]uuid.UUID, error)

	AddPermissionToDeletingList(ctx context.Context, key string) error
	RemovePermissionFromDeletingList(ctx context.Context, key string) error
	GetRandomPermissionFromDeletingList(ctx context.Context) (perm string, isListEmpty bool, err error)

	GetUserPermissions(ctx context.Context, userUid uuid.UUID) ([]entity.UserPermission, error)
	UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error
}

type UsersRepository interface {
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error)
}
