package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) Logout(ctx context.Context, device entity.UserDevice) error {
	log.Printf("usecaseAuth.Logout: user %s", device.UserUid)

	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		if uuid.Equal(device.UserUid, uuid.UUID{}) {
			return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
		}

		if err := uc.ucDevices.HandleLogoutFromDevice(ctx, device, time.Now()); err != nil {
			return err
		}

		if err := uc.sessionsRepo.DeleteSessionByUserUidAndDeviceUid(ctx, device.UserUid, device.Uid); err != nil {
			log.Printf("usecaseAuth.Logout: failed to delete session for user %s device %s: %v", device.UserUid, device.Uid, err)
			return err
		}
		return nil
	})
}
