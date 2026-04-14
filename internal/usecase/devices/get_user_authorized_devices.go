package ucDevices

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCase) GetUserAuthorizedDevices(ctx context.Context, userUid uuid.UUID) ([]entity.UserAuthorizedDevice, error) {
	return uc.devicesRepo.GetUserAuthorizedDevices(ctx, userUid)
}
