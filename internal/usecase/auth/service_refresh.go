package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
)

func (uc *UseCaseAuth) ServiceRefresh(ctx context.Context, token string) (string, string, error) {
	log.Printf("usecaseAuth.ServiceRefresh: ")
	refreshTime := time.Now()

	var access, refresh string
	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		tokenInfo, err := uc.parseAndVerifyRefreshToken(ctx, token)
		if err != nil {
			return err
		}

		service, err := uc.getServiceWithRolesByUid(ctx, tokenInfo.ServiceUid)
		if err != nil {
			return err
		}

		newTokenInfo := entity.NewSystemTokenInfo(service, tokenInfo.DeviceUid, tokenInfo.SessionUid, refreshTime)
		access, refresh, err = uc.jwtManager.NewSystemTokens(newTokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
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
