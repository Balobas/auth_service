package ucDevices

import (
	"context"
	"time"

	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCase) UnauthorizeUserDevices(ctx context.Context, userUid uuid.UUID, unauthTime time.Time) error {
	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		authDevices, err := uc.devicesRepo.GetUserAuthorizedDevices(ctx, userUid)
		if err != nil {
			return err
		}

		for _, device := range authDevices {
			device = device.Unauthorize(unauthTime)

			if err := uc.devicesRepo.UpdateAuthorizedDevice(ctx, device); err != nil {
				return err
			}
		}
		return nil
	})
}
