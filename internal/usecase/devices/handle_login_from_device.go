package ucDevices

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
)

func (uc *UseCase) HandleLoginFromDevice(ctx context.Context, device entity.Device, loginTime time.Time) error {
	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		authDevice, isFound, err := uc.devicesRepo.GetAuthorizedDevice(ctx, device.MaintainerUid, device.Uid)
		if err != nil {
			return err
		}
		if isFound {
			// Истек токен, владелец сам не разлогинивался
			if authDevice.UnauthorizedAt == nil {
				return nil
			}

			// иначе, владелец разлогинивался с девайса и сейчас нужно учесть его логин
			authDevice = authDevice.Authorize(loginTime)
			return uc.devicesRepo.UpdateAuthorizedDevice(ctx, authDevice)
		}

		authDevice = device.Authorize(loginTime)
		return uc.devicesRepo.CreateAuthorizedDevice(ctx, authDevice)
	})
}
