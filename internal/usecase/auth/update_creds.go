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

func (uc *UseCaseAuth) UpdateUserCreds(ctx context.Context, user entity.User, password string, deviceUid uuid.UUID) (string, string, error) {
	log.Printf("usecaseAuth.UpdateUserCreds: user %s", user.Uid)
	now := time.Now()

	if uuid.Equal(user.Uid, uuid.UUID{}) {
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid"))
	}
	if uuid.Equal(deviceUid, uuid.UUID{}) {
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty device uid"))
	}
	if len(user.Email) == 0 && len(password) == 0 {
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user email and password"))
	}

	var access, refresh string

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		// TODO: возвращать юзера, иначе может быть пустой емэйл в токене
		if err := uc.ucUsers.UpdateUser(ctx, user, password); err != nil {
			return err
		}

		roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
		if err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to get user %s roles: %v", user.Uid, err)
			return err
		}

		rolesStrs := entity.RolesToStrings(roles)
		user.Roles = rolesStrs

		if err := uc.ucDevices.UnauthorizeUserDevices(ctx, user.Uid, now); err != nil {
			return err
		}
		if err := uc.sessionsRepo.DeleteUserSessions(ctx, user.Uid); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to delete all sessions for user %s: %v", user.Uid, err)
			return err
		}

		if err := uc.ucDevices.HandleLoginFromDevice(
			ctx,
			entity.Device{Uid: deviceUid, Type: entity.DeviceTypeUserDevice, MaintainerUid: user.Uid},
			now,
		); err != nil {
			return err
		}

		session := entity.NewUserSession(uuid.NewV4(), user.Uid, deviceUid, now)
		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to create new session for user %s: %v", user.Uid, err)
			return err
		}

		tokenInfo := entity.NewUserTokenInfo(user, deviceUid, session.Uid, now)
		access, refresh, err = uc.jwtManager.NewUserTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
		return err
	}); err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
