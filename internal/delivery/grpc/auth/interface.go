package deliveryGrpcAuth

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type UcUsers interface {
	Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error)
	GetAdmins(ctx context.Context) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, password string) error
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	DeleteUser(ctx context.Context, userUid uuid.UUID) error
}

type UcAuth interface {
	Login(ctx context.Context, params entity.LoginParams) (string, string, error)
	Logout(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) error
	Refresh(ctx context.Context, token string) (string, string, error)
	VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error)
	VerifyAccess(ctx context.Context, uri string, method string, token string) error
	UpdateUserCreds(ctx context.Context, user entity.User, password string, deviceUid uuid.UUID) (string, string, error)
}

type UcVerification interface {
	Verify(ctx context.Context, token string) error
}

type UcDevices interface {
	GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.AuthorizedDevice, error)
	GetUsersAuthorizedDevices(ctx context.Context, usersUids ...uuid.UUID) ([]entity.AuthorizedDevice, error)
}

type UcAccess interface {
	CreatePermission(ctx context.Context, perm entity.Permission) error
	DeletePermission(ctx context.Context, key string) error
	GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error)
	UpdatePermission(ctx context.Context, perm entity.Permission) error

	CreateRole(ctx context.Context, role entity.Role) error
	DeleteRole(ctx context.Context, role string) error
	GetRole(ctx context.Context, role string) (entity.Role, error)
	UpdateRole(ctx context.Context, role entity.Role) error
	GetRoles(ctx context.Context, rolePattern string, limit int64, offset int64) ([]entity.Role, error)

	AddPermissionsToRole(ctx context.Context, role string, permissions []string) error
	DeletePermissionsFromRole(ctx context.Context, role string, permissions []string) error
	GetRolePermissions(ctx context.Context, role string) ([]entity.Permission, error)

	AddRoleToUser(ctx context.Context, userUid uuid.UUID, role string) error
	RemoveRoleFromUser(ctx context.Context, userUid uuid.UUID, role string) error
	GetUserRoles(ctx context.Context, userUid uuid.UUID) ([]entity.Role, error)

	AddResourcePermission(ctx context.Context, uri string, method string, permission string) error
	DeleteResourcePermission(ctx context.Context, uri string, method string, permission string) error
	GetResourcePermission(ctx context.Context, uri string, method string) (string, error)
	UpdateResourcePermission(ctx context.Context, uri string, method string, permission string) error
	GetResourcesPermissions(ctx context.Context, limit int64, offset int64) ([]entity.ResourcePermissions, error)
}
