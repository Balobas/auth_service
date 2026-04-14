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

func (uc *UseCaseAuth) UpdateUserCreds(ctx context.Context, user entity.User, password string, device entity.UserDevice) (string, string, error) {
	log.Printf("usecaseAuth.UpdateUserCreds: user %s", user.Uid)

	if uuid.Equal(user.Uid, uuid.UUID{}) {
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid"))
	}
	if len(user.Email) == 0 && len(password) == 0 {
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user email and password"))
	}
	if err := device.Validate(); err != nil {
		return emptyTokensWithError(err)
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

		now := time.Now()

		session := entity.Session{
			Uid:            uuid.NewV4(),
			UserUid:        user.Uid,
			DeviceUid:      device.Uid,
			TokensIssuedAt: now.Unix(),
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		tokenInfo := entity.TokenInfo{
			UserUid:    user.Uid,
			DeviceUid:  device.Uid,
			Email:      user.Email,
			Roles:      rolesStrs,
			SessionUid: session.Uid,
			IssuedAt:   now.Unix(),
		}

		if err := uc.ucDevices.UnauthorizeUserDevices(ctx, user.Uid, now); err != nil {
			return err
		}
		if err := uc.ucDevices.HandleLoginFromDevice(ctx, device, now); err != nil {
			return err
		}

		if err := uc.sessionsRepo.DeleteSessionsByUserUid(ctx, user.Uid); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to delete all sessions for user %s: %v", user.Uid, err)
			return err
		}

		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to create new session for user %s: %v", user.Uid, err)
			return err
		}

		access, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
		if err != nil {
			return errors.Wrapf(err, "failed to build jwt")
		}
		refresh, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
		if err != nil {
			return errors.Wrapf(err, "failed to build jwt")
		}

		return nil
	}); err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
