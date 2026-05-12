package ucDevices

import (
	"context"
	"errors"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCase) GetUsersAuthorizedDevices(ctx context.Context, usersUids ...uuid.UUID) ([]entity.AuthorizedDevice, error) {
	if len(usersUids) == 0 {
		return nil, errors.New("empty users uids")
	}

	return uc.devicesRepo.GetUsersAuthorizedDevices(ctx, usersUids...)
}
