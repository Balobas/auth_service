package ucDevices

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type (
	DevicesRepository interface {
		CreateAuthorizedDevice(ctx context.Context, device entity.AuthorizedDevice) error
		UpdateAuthorizedDevice(ctx context.Context, device entity.AuthorizedDevice) error
		GetAuthorizedDevice(ctx context.Context, maintainerUid uuid.UUID, deviceUid uuid.UUID) (entity.AuthorizedDevice, bool, error)
		GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.AuthorizedDevice, error)
		GetUsersAuthorizedDevices(ctx context.Context, usersUids ...uuid.UUID) ([]entity.AuthorizedDevice, error)
	}
)
