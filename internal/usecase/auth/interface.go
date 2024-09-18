package useCaseAuth

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type Config interface {
	MinPasswordLen() int
	AccessJwtTTL() time.Duration
	RefreshJwtTTL() time.Duration
}

type UcUsers interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, bool, error)
	UpdateUser(ctx context.Context, user entity.User, password string) error
}

type UcCredentials interface {
	Validate(ctx context.Context, userUid uuid.UUID, password string) error
}

type JwtManager interface {
	NewToken(info entity.TokenInfo, ttl time.Duration) (string, error)
	ParseToken(tokenStr string) (entity.TokenInfo, error)
}

type SessionsRepository interface {
	CreateSession(ctx context.Context, session entity.Session) error
	GetSessionByUid(ctx context.Context, uid uuid.UUID) (entity.Session, bool, error)
	GetSessionByUserUid(ctx context.Context, userUid uuid.UUID) (entity.Session, bool, error)
	UpdateSession(ctx context.Context, sessionUid uuid.UUID, updatedAt time.Time) error
	DeleteSessionByUid(ctx context.Context, uid uuid.UUID) error
	DeleteSessionByUserUid(ctx context.Context, userUid uuid.UUID) error
}

type PermissionsRepository interface {
	GetUserPermissions(ctx context.Context, userUid uuid.UUID) ([]entity.UserPermission, error)
}
