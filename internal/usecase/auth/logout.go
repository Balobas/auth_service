package useCaseAuth

import (
	"context"
	"log"
	"time"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Logout выполняет логаут пользователя с одного из его девайсов
// Если передан девайс, не принадлежащий юзеру - вернется ошибка
// Допустимо разлогинивать только свои девайсы
func (uc *UseCaseAuth) Logout(ctx context.Context, userUid uuid.UUID, deviceUid uuid.UUID) error {
	log.Printf("usecaseAuth.Logout: user %s, device %s", userUid, deviceUid)
	if uuid.Equal(deviceUid, uuid.UUID{}) {
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty device uid")
	}

	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		if err := uc.ucDevices.HandleLogoutFromDevice(ctx, userUid, deviceUid, time.Now()); err != nil {
			return err
		}

		if err := uc.sessionsRepo.DeleteSessionByUserUidAndDeviceUid(ctx, userUid, deviceUid); err != nil {
			log.Printf("usecaseAuth.Logout: failed to delete session for user %s device %s: %v", userUid, deviceUid, err)
			return err
		}
		return nil
	})
}
