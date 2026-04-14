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
	log.Printf("usecaseAuth.Login: email %s", params.Email)

	if err := validations.ValidateEmail(params.Email); err != nil {
		return emptyTokensWithError(err)
	}

	if err := params.Device.Validate(); err != nil {
		return emptyTokensWithError(err)
	}

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		user, err := uc.getUserWithRolesByEmail(ctx, params.Email)
		if err != nil {
			return err
		}

		if err := uc.ucCredentials.Validate(ctx, user.Uid, params.Password); err != nil {
			return err
		}

		savedSession, isFound, err := uc.sessionsRepo.GetSessionByUserUidAndDeviceUid(ctx, user.Uid, params.Device.Uid)
		if err != nil {
			return err
		}
		if isFound {
			access, refresh, err = uc.handleLoginWithExistingSession(user, savedSession)
			return err
		}

		loginTime := time.Now()

		if err := uc.ucDevices.HandleLoginFromDevice(
			ctx,
			params.Device.WithUserUid(user.Uid),
			loginTime,
		); err != nil {
			return err
		}

		session := entity.Session{
			Uid:            uuid.NewV4(),
			UserUid:        user.Uid,
			DeviceUid:      params.Device.Uid,
			CreatedAt:      loginTime,
			TokensIssuedAt: loginTime.Unix(),
		}

		tokenInfo := entity.TokenInfo{
			UserUid:    user.Uid,
			DeviceUid:  params.Device.Uid,
			Email:      user.Email,
			Roles:      user.Roles,
			SessionUid: session.Uid,
			IssuedAt:   loginTime.Unix(),
		}

		access, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
		if err != nil {
			return err
		}
		refresh, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
		if err != nil {
			return err
		}

		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return emptyTokensWithError(err)
	}

	return access, refresh, nil
}

func emptyTokensWithError(err error) (string, string, error) {
	return "", "", err
}

func (uc *UseCaseAuth) handleLoginWithExistingSession(user entity.User, session entity.Session) (string, string, error) {
	tokenInfo := entity.TokenInfo{
		UserUid:    session.UserUid,
		DeviceUid:  session.DeviceUid,
		Email:      user.Email,
		Roles:      user.Roles,
		SessionUid: session.Uid,
		IssuedAt:   session.TokensIssuedAt,
	}

	access, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		return emptyTokensWithError(err)
	}
	refresh, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		return emptyTokensWithError(err)
	}
	return access, refresh, nil
}
