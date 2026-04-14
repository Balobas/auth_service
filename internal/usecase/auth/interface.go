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
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, password string) error
}

type UcCredentials interface {
	Validate(ctx context.Context, userUid uuid.UUID, password string) error
}

type UcDevices interface {
	HandleLoginFromDevice(ctx context.Context, device entity.UserDevice, loginTime time.Time) error
	HandleLogoutFromDevice(ctx context.Context, device entity.UserDevice, logoutTime time.Time) error
	UnauthorizeUserDevices(ctx context.Context, userUid uuid.UUID, unauthTime time.Time) error
}

type JwtManager interface {
	NewToken(info entity.TokenInfo, ttl time.Duration) (string, error)
	ParseToken(tokenStr string) (entity.TokenInfo, error)
}

type SessionsRepository interface {
	CreateSession(ctx context.Context, session entity.Session) error
	GetSessionByUid(ctx context.Context, uid uuid.UUID) (entity.Session, bool, error)
	GetSessionByUserUidAndDeviceUid(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) (entity.Session, bool, error)
	UpdateSession(ctx context.Context, session entity.Session) error
	DeleteSessionsByUserUid(ctx context.Context, userUid uuid.UUID) error
	DeleteSessionByUserUidAndDeviceUid(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) error
}

type AccessRepository interface {
	GetUserRoles(ctx context.Context, userUid uuid.UUID) ([]entity.Role, error)
	IsUserHasPermissionsForResource(ctx context.Context, userUid uuid.UUID, uri string, method string) (bool, error)
}
