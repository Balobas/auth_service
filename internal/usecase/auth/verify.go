package useCaseAuth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error) {
	log.Printf("auth.VerifyAuth")

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("verifyAuth: failed to parse token\n")
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("token expired")
		return entity.TokenInfo{}, errors.New("token expired")
	}

	user, isFound, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("user not found: %v\n", err)
		return entity.TokenInfo{}, errors.New("user not found")
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("user in token is not user in request")
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("failed to get permissions: %v", err)
		return entity.TokenInfo{}, errors.Wrap(err, "failed to get user permissions")
	}
	user.Permissions = perms

	_, isFound, err = uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("failed to get session: %v", err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("session not found")
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	if len(tokenInfo.Permissions) != len(user.Permissions) {
		log.Printf("invalid permissions in token")
		return tokenInfo, ErrPermissionsNotMatch
	}

	permsMap := make(map[entity.UserPermission]struct{}, len(tokenInfo.Permissions))
	for _, perm := range tokenInfo.Permissions {
		permsMap[entity.UserPermission(perm)] = struct{}{}
	}

	for _, perm := range user.Permissions {
		if _, ok := permsMap[perm]; !ok {
			log.Printf("invalid permissions in token")
			return tokenInfo, ErrPermissionsNotMatch
		}
	}

	if tokenInfo.Role != string(user.Role) {
		log.Printf("role from token is invalid: token role '%s' user role: '%s'", tokenInfo.Role, user.Role)
		return tokenInfo, ErrRoleNotMatch
	}

	return tokenInfo, nil
}

var (
	ErrRoleNotMatch        = fmt.Errorf("role not match")
	ErrPermissionsNotMatch = fmt.Errorf("permissions not match")
)
