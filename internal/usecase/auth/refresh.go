package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) Refresh(ctx context.Context, token string) (string, string, error) {
	log.Printf("usecaseAuth.Refresh: ")
	refreshTime := time.Now()

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		tokenInfo, err := uc.parseAndVerifyRefreshToken(ctx, token)
		if err != nil {
			return err
		}

		user, err := uc.getUserWithRolesByUid(ctx, tokenInfo.UserUid)
		if err != nil {
			return err
		}

		newTokenInfo := entity.NewUserTokenInfo(user, tokenInfo.DeviceUid, tokenInfo.SessionUid, refreshTime)
		access, refresh, err = uc.jwtManager.NewUserTokens(newTokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
		if err != nil {
			return err
		}

		if err := uc.sessionsRepo.UpdateSession(ctx, entity.Session{
			Uid:            tokenInfo.SessionUid,
			MaintainerUid:  tokenInfo.UserUid,
			DeviceUid:      tokenInfo.DeviceUid,
			TokensIssuedAt: refreshTime.Unix(),
			UpdatedAt:      refreshTime,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return emptyTokensWithError(err)
	}

	return access, refresh, nil
}

func (uc *UseCaseAuth) getUserWithRolesByUid(ctx context.Context, uid uuid.UUID) (entity.User, error) {
	user, err := uc.ucUsers.GetUserByUid(ctx, uid)
	if err != nil {
		return entity.User{}, err
	}
	return uc.withRoles(ctx, user)
}

func (uc *UseCaseAuth) getUserWithRolesByEmail(ctx context.Context, email string) (entity.User, error) {
	user, err := uc.ucUsers.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, err
	}
	return uc.withRoles(ctx, user)
}

func (uc *UseCaseAuth) withRoles(ctx context.Context, user entity.User) (entity.User, error) {
	roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
	if err != nil {
		return entity.User{}, err
	}
	user.Roles = entity.RolesToStrings(roles)
	return user, nil
}
