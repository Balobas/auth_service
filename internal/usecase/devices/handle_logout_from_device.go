package ucDevices

import (
	"context"
	"fmt"
	"time"

	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCase) HandleLogoutFromDevice(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID, logoutTime time.Time) error {
	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		authDevice, isFound, err := uc.devicesRepo.GetUserAuthorizedDevice(ctx, userUid, deviceUid)
		if err != nil {
			return err
		}

		if !isFound {
			// нет связи юзера с девайсом,
			// девайс должен быть учтен, это аномальный случай когда данные на вход валидные, а девайса нет
			// в остальных случаях значит, что юзер пытается разлогинить другой девайс - ошибка
			return fmt.Errorf("invalid device")
		}

		if authDevice.UnauthorizedAt != nil {
			return fmt.Errorf("device already logout")
		}

		authDevice = authDevice.Unauthorize(logoutTime)
		return uc.devicesRepo.UpdateUserAuthorizedDevice(ctx, authDevice)
	})
}
