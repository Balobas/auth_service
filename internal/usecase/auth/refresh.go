package useCaseAuth

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func (uc *UseCaseAuth) Refresh(ctx context.Context, token string) (string, string, error) {
	tokenInfo, err := uc.VerifyAuth(ctx, token)
	if err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	access, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}
	refresh, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}

	if err := uc.sessionsRepo.UpdateSession(ctx, tokenInfo.SessionUid, time.Now()); err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
