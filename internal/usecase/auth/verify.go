package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error) {
	log.Printf("usecaseAuth.VerifyAuth: token: %v", token)

	if len(token) == 0 {
		log.Printf("usecaseAuth.VerifyAuth: empty token")
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrBadRequest, "empty token")
	}

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to parse token: %v\n", err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("usecaseAuth.VerifyAuth: token expired for user %s", tokenInfo.UserUid)
		return entity.TokenInfo{}, serviceErrors.ErrTokenExpired
	}

	user, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get user %s: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("usecaseAuth.VerifyAuth: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "user in token is not user in request")
	}

	perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get user %s permissions: %v", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.Wrap(err, "failed to get user permissions")
	}
	user.Permissions = perms

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseAuth.VerifyAuth: session %s not found", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrNotFound, "session")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("usecaseAuth.VerifyAuth: token from request already invalid for session %s", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "you should use last issued token")
	}

	if len(tokenInfo.Permissions) != len(user.Permissions) {
		log.Printf("usecaseAuth.VerifyAuth: invalid permissions in user %s token", tokenInfo.UserUid)
		return tokenInfo, errors.Wrap(serviceErrors.ErrInvalidToken, "invalid permissions in token")
	}

	permsMap := make(map[entity.UserPermission]struct{}, len(tokenInfo.Permissions))
	for _, perm := range tokenInfo.Permissions {
		permsMap[entity.UserPermission(perm)] = struct{}{}
	}

	for _, perm := range user.Permissions {
		if _, ok := permsMap[perm]; !ok {
			log.Printf("usecaseAuth.VerifyAuth: invalid permissions in user %s token", tokenInfo.UserUid)
			return tokenInfo, errors.Wrap(serviceErrors.ErrInvalidToken, "invalid permissions in token")
		}
	}

	if tokenInfo.Role != string(user.Role) {
		log.Printf("usecaseAuth.VerifyAuth: role from token is invalid: token role '%s' user role: '%s' user uid '%s'", tokenInfo.Role, user.Role, user.Uid)
		return tokenInfo, errors.Wrap(serviceErrors.ErrInvalidToken, "invalid user role in token")
	}

	return tokenInfo, nil
}
