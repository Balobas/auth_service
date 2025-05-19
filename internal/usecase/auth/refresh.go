package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCaseAuth) Refresh(ctx context.Context, token string) (string, string, error) {
	log.Printf("usecaseAuth.Refresh: ")
	tokenInfo, err := uc.verifyRefreshToken(ctx, token)
	if err != nil {
		log.Printf("usecaseAuth.Refresh: failed to validate token: %v", err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	user, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("usecaseAuth.Refresh: failed to get user by email %s: %v", tokenInfo.Email, err)
		return emptyTokensWithError(errors.WithStack(err))
	}
	perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseAuth.Refresh: failed to get user %s permissions: %v", user.Uid, err)
		return emptyTokensWithError(errors.WithStack(err))
	}
	user.Permissions = perms

	refreshTime := time.Now()

	newTokenInfo := entity.TokenInfo{
		UserUid:     user.Uid,
		Email:       user.Email,
		Permissions: user.PermissionsStrings(),
		Role:        string(user.Role),
		SessionUid:  tokenInfo.SessionUid,
		IssuedAt:    refreshTime.Unix(),
	}

	access, err := uc.jwtManager.NewToken(newTokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		log.Printf("usecaseAuth.Refresh: failed to build jwt token for user %s: %v", user.Uid, err)
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}
	refresh, err := uc.jwtManager.NewToken(newTokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		log.Printf("usecaseAuth.Refresh: failed to build jwt token for user %s: %v", user.Uid, err)
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}

	if err := uc.sessionsRepo.UpdateSession(ctx, entity.Session{
		Uid:            tokenInfo.SessionUid,
		TokensIssuedAt: refreshTime.Unix(),
		UpdatedAt:      time.Now(),
	}); err != nil {
		log.Printf("usecaseAuth.Refresh: failed to update session for user %s: %v", user.Uid, err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
