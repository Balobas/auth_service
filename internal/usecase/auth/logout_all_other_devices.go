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

// LogoutAllOtherDevices
// Выполняет логаут всех устройств юзера кроме текущего
func (uc *UseCaseAuth) LogoutAllOtherDevices(ctx context.Context, userUid uuid.UUID, keepAuthorizedDeviceUid uuid.UUID) error {
	log.Printf("usecaseAuth.LogoutAllOtherDevices: user %s, keepAuthorizedDeviceUid %s", userUid, keepAuthorizedDeviceUid)

	if uuid.Equal(keepAuthorizedDeviceUid, uuid.UUID{}) {
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty device uid")
	}
	now := time.Now()

	return uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		keepActiveSession, isFound, err := uc.sessionsRepo.GetSessionByMaintainerAndDevice(ctx, userUid, keepAuthorizedDeviceUid)
		if err != nil {
			return err
		}
		if !isFound {
			return errors.Wrap(serviceErrors.ErrNotFound, "session")
		}

		if err := uc.ucDevices.UnauthorizeUserDevices(ctx, userUid, now, keepAuthorizedDeviceUid); err != nil {
			return err
		}

		if err := uc.sessionsRepo.DeleteUserSessions(ctx, userUid, keepActiveSession.Uid); err != nil {
			log.Printf("usecaseAuth.LogoutAllOtherDevices: failed to delete session for user %s: %v", userUid, err)
			return err
		}
		return nil
	})
}
