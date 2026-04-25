package ucDevices

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
)

func (uc *UseCase) HandleLoginFromDevice(ctx context.Context, device entity.UserDevice, loginTime time.Time) error {
	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		authDevice, isFound, err := uc.devicesRepo.GetUserAuthorizedDevice(ctx, device.UserUid, device.Uid)
		if err != nil {
			return err
		}
		if isFound {
			// Истек токен, юзер сам не разлогинивался
			if authDevice.UnauthorizedAt == nil {
				return nil
			}

			// иначе, юзер разлогинивался с девайса и сейчас нужно учесть его логин
			authDevice = authDevice.Authorize(loginTime)
			return uc.devicesRepo.UpdateUserAuthorizedDevice(ctx, authDevice)
		}

		authDevice = device.Authorize(loginTime)
		return uc.devicesRepo.CreateUserAuthorizedDevice(ctx, authDevice)
	})
}
