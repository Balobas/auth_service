package useCaseUsers

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
	GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error)
	GetByEmail(ctx context.Context, email string) (entity.User, bool, error)
	UpdateUser(ctx context.Context, userParams entity.User) error
	DeleteUser(ctx context.Context, uid uuid.UUID) error
}

type PermissionsRepository interface {
	CreateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error
	UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error
}

type UcVerification interface {
	CreateVerification(ctx context.Context, userUid uuid.UUID, email string) error
}

type UcCredentials interface {
	Validate(ctx context.Context, userUid uuid.UUID, password string) error
	Update(ctx context.Context, userUid uuid.UUID, password string) error
	Create(ctx context.Context, userUid uuid.UUID, password string) error
}
