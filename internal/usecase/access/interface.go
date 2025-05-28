package useCaseAccess

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type AccessRepository interface {
	CreatePermission(ctx context.Context, perm entity.Permission) error
	UpdatePermission(ctx context.Context, perm entity.Permission) error
	GetPermission(ctx context.Context, key string) (entity.Permission, bool, error)
	GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error)
	DeletePermission(ctx context.Context, key string) error
	RemovePermissionFromUsers(ctx context.Context, key string, limitUsers int64) ([]uuid.UUID, error)

	CreateRole(ctx context.Context, role entity.Role) error
	UpdateRole(ctx context.Context, role entity.Role) error
	DeleteRole(ctx context.Context, role string) error
	GetRole(ctx context.Context, role string) (entity.Role, bool, error)
	GetRoles(ctx context.Context, rolePattern string, limit int64, offset int64) ([]entity.Role, error)

	AddRoleToUser(ctx context.Context, userUid uuid.UUID, role string) error
	DeleteRoleFromUser(ctx context.Context, userUid uuid.UUID, role string) error
	GetUserRoles(ctx context.Context, userUid uuid.UUID) ([]entity.Role, error)

	AddPermissionsToRole(ctx context.Context, role string, permissions []string) error
	DeletePermissionsFromRole(ctx context.Context, role string, permissions []string) error
	GetRolePermissions(ctx context.Context, role string) ([]entity.Permission, error)

	AddResourcePermission(ctx context.Context, uri string, method string, permission string) error
	GetResourcePermission(ctx context.Context, uri string, method string) (string, bool, error)
	UpdateResourcePermission(ctx context.Context, uri string, method string, permission string) error
	DeleteResourcePermission(ctx context.Context, uri string, method string, permission string) error
	GetResourcesPermissions(ctx context.Context, limit int64, offset int64) ([]entity.ResourcePermissions, error)
}

type UsersRepository interface {
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error)
}
