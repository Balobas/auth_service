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
		if errors.Is(err, serviceErrors.ErrNotFound) {
			log.Printf("usecaseAuth.VerifyAuth: user with email %s not found", tokenInfo.Email)
			return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "user not found")
		}

		log.Printf("usecaseAuth.VerifyAuth: failed to get user %s: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("usecaseAuth.VerifyAuth: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "user in token is not user in request")
	}

	roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get user %s roles: %v", user.Uid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	user.Roles = entity.RolesToStrings(roles)

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseAuth.VerifyAuth: session %s not found", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "session not found")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("usecaseAuth.VerifyAuth: token from request already invalid for session %s", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrInvalidToken, "you should use last issued token")
	}

	if len(tokenInfo.Roles) != len(user.Roles) {
		log.Printf("usecaseAuth.VerifyAuth: invalid roles in user %s token: len mismatch", tokenInfo.UserUid)
		return tokenInfo, errors.Wrap(serviceErrors.ErrInvalidToken, "invalid roles in token")
	}

	rolesMap := make(map[string]struct{}, len(tokenInfo.Roles))
	for _, role := range tokenInfo.Roles {
		rolesMap[role] = struct{}{}
	}

	for _, role := range user.Roles {
		if _, ok := rolesMap[role]; !ok {
			log.Printf("usecaseAuth.VerifyAuth: invalid roles in user %s token: role %s doesnt exists in user", tokenInfo.UserUid, role)
			return tokenInfo, errors.Wrap(serviceErrors.ErrInvalidToken, "invalid roles in token")
		}
	}

	return tokenInfo, nil
}
