package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

/*
Если юзер уже авторизован с девайса, то новая сессия не создастся
вернутся токены существующей сессии

Для одного юзера и девайса в один момент времени может быть только одна сессия
*/
func (uc *UseCaseAuth) Login(ctx context.Context, params entity.LoginParams) (string, string, error) {
	log.Printf("usecaseAuth.Login: email %s, device %s", params.Email, params.Device.Uid)

	if err := validations.ValidateEmail(params.Email); err != nil {
		return emptyTokensWithError(err)
	}
	loginTime := time.Now()

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		user, err := uc.getUserWithRolesByEmail(ctx, params.Email)
		if err != nil {
			return err
		}

		device, err := entity.NewUserDevice(user.Uid, params.Device)
		if err != nil {
			return err
		}

		if err := uc.ucCredentials.Validate(ctx, user.Uid, params.Password); err != nil {
			return err
		}

		savedSession, isFound, err := uc.sessionsRepo.GetSessionByMaintainerAndDevice(ctx, user.Uid, device.Uid)
		if err != nil {
			return err
		}
		if isFound {
			access, refresh, err = uc.handleUserLoginWithExistingSession(user, savedSession)
			return err
		}

		if err := uc.ucDevices.HandleLoginFromDevice(ctx, device, loginTime); err != nil {
			return err
		}

		session := entity.NewUserSession(uuid.NewV4(), user.Uid, device.Uid, loginTime)
		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			return err
		}

		tokenInfo := entity.NewUserTokenInfo(user, device.Uid, session.Uid, loginTime)
		access, refresh, err = uc.jwtManager.NewUserTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
		return err
	}); err != nil {
		return emptyTokensWithError(err)
	}

	return access, refresh, nil
}

func emptyTokensWithError(err error) (string, string, error) {
	return "", "", err
}

func (uc *UseCaseAuth) handleUserLoginWithExistingSession(user entity.User, session entity.Session) (string, string, error) {
	tokenInfo := entity.UserTokenInfoFromSession(user, session)
	return uc.jwtManager.NewUserTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
}
