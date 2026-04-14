package ucDevices

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
)

func (uc *UseCase) HandleLogoutFromDevice(ctx context.Context, device entity.UserDevice, logoutTime time.Time) error {
	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		authDevice, isFound, err := uc.devicesRepo.GetUserAuthorizedDevice(ctx, device.UserUid, device.Uid)
		if err != nil {
			return err
		}

		if !isFound {
			// критическая ситуация, девайс не был сохранен
			// вернуть ошибку нельзя, нужно просто сохранить девайс с временем логаута
			authDevice = device.Authorize(logoutTime).Unauthorize(logoutTime)
			return uc.devicesRepo.CreateUserAuthorizedDevice(ctx, authDevice)
		}

		authDevice = authDevice.Unauthorize(logoutTime)
		return uc.devicesRepo.UpdateUserAuthorizedDevice(ctx, authDevice)
	})
}
