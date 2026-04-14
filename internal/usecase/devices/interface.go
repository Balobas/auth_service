package ucDevices

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type (
	DevicesRepository interface {
		CreateUserAuthorizedDevice(ctx context.Context, device entity.UserAuthorizedDevice) error
		UpdateUserAuthorizedDevice(ctx context.Context, device entity.UserAuthorizedDevice) error
		GetUserAuthorizedDevice(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) (entity.UserAuthorizedDevice, bool, error)
		GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.UserAuthorizedDevice, error)
	}
)
