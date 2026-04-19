package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
)

func (uc *UseCaseAuth) Refresh(ctx context.Context, token string) (string, string, error) {
	log.Printf("usecaseAuth.Refresh: ")

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		tokenInfo, err := uc.parseAndVerifyRefreshToken(ctx, token)
		if err != nil {
			return err
		}

		user, err := uc.getUserWithRolesByEmail(ctx, tokenInfo.Email)
		if err != nil {
			return err
		}

		refreshTime := time.Now()

		newTokenInfo := entity.TokenInfo{
			UserUid:    user.Uid,
			DeviceUid:  tokenInfo.DeviceUid,
			Email:      user.Email,
			Roles:      user.Roles,
			SessionUid: tokenInfo.SessionUid,
			IssuedAt:   refreshTime.Unix(),
		}

		access, err = uc.jwtManager.NewToken(newTokenInfo, uc.cfg.AccessJwtTTL())
		if err != nil {
			return err
		}
		refresh, err = uc.jwtManager.NewToken(newTokenInfo, uc.cfg.RefreshJwtTTL())
		if err != nil {
			return err
		}

		if err := uc.sessionsRepo.UpdateSession(ctx, entity.Session{
			Uid:            tokenInfo.SessionUid,
			UserUid:        tokenInfo.UserUid,
			DeviceUid:      tokenInfo.DeviceUid,
			TokensIssuedAt: refreshTime.Unix(),
			UpdatedAt:      time.Now(),
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return emptyTokensWithError(err)
	}

	return access, refresh, nil
}

func (uc *UseCaseAuth) getUserWithRolesByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := uc.ucUsers.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}
	roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
	if err != nil {
		return entity.User{}, err
	}

	rolesStrs := entity.RolesToStrings(roles)
	user.Roles = rolesStrs
	return user, nil
}
