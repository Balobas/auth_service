package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCaseAuth) Refresh(ctx context.Context, token string) (string, string, error) {
	tokenInfo, err := uc.VerifyAuth(ctx, token)
	switch {
	case err == nil:
	case errors.Is(err, ErrRoleNotMatch) || errors.Is(err, ErrPermissionsNotMatch):
	default:
		return emptyTokensWithError(errors.WithStack(err))
	}

	user, _, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("failed to get user")
		return emptyTokensWithError(errors.WithStack(err))
	}
	perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("failed to get permissions")
		return emptyTokensWithError(err)
	}
	user.Permissions = perms

	newTokenInfo := entity.TokenInfo{
		UserUid:     user.Uid,
		Email:       user.Email,
		Permissions: user.PermissionsStrings(),
		Role:        string(user.Role),
		SessionUid:  tokenInfo.SessionUid,
	}

	access, err := uc.jwtManager.NewToken(newTokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}
	refresh, err := uc.jwtManager.NewToken(newTokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}

	if err := uc.sessionsRepo.UpdateSession(ctx, tokenInfo.SessionUid, time.Now()); err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
