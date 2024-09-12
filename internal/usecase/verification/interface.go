package useCaseVerification

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type Config interface {
	VerificationTokenLen() int64
}

type VerificationRepository interface {
	CreateVerification(ctx context.Context, verification entity.Verification) error
	GetUserVerification(ctx context.Context, userUid uuid.UUID) (entity.Verification, error)
	GetVerificationByToken(ctx context.Context, token string) (entity.Verification, bool, error)
	UpdateVerification(ctx context.Context, verification entity.Verification) error
	DeleteVerification(ctx context.Context, userUid uuid.UUID) error
}

type PermissionsRepository interface {
	UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error
}
