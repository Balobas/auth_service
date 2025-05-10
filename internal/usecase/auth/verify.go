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
	log.Printf("auth.VerifyAuth: token: %v", token)

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("auth.VerifyAuth: failed to parse token: %v\n", err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("auth.VerifyAuth: token expired for user %s", tokenInfo.UserUid)
		return entity.TokenInfo{}, errors.New("token expired")
	}

	user, isFound, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("auth.VerifyAuth: failed to get user %s: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("auth.VerifyAuth: user %s not found: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.New("user not found")
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("auth.VerifyAuth: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("auth.VerifyAuth: failed to get user %s permissions: %v", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.Wrap(err, "failed to get user permissions")
	}
	user.Permissions = perms

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("auth.VerifyAuth: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("auth.VerifyAuth: session %s not found", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		fmt.Println("in session: ", session.TokensIssuedAt)
		fmt.Println("in token: ", tokenInfo.IssuedAt)
		log.Printf("auth.VerifyAuth: token from request already invalid for session %s", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.New("invalid token. you should use last issued token")
	}

	if len(tokenInfo.Permissions) != len(user.Permissions) {
		log.Printf("auth.VerifyAuth: invalid permissions in user %s token", tokenInfo.UserUid)
		return tokenInfo, ErrPermissionsNotMatch
	}
	log.Printf("token perms: %v", tokenInfo.Permissions)
	log.Printf("user permissions: %v", user.Permissions)

	permsMap := make(map[entity.UserPermission]struct{}, len(tokenInfo.Permissions))
	for _, perm := range tokenInfo.Permissions {
		permsMap[entity.UserPermission(perm)] = struct{}{}
	}

	for _, perm := range user.Permissions {
		if _, ok := permsMap[perm]; !ok {
			log.Printf("auth.VerifyAuth: invalid permissions in user %s token", tokenInfo.UserUid)
			return tokenInfo, ErrPermissionsNotMatch
		}
	}

	if tokenInfo.Role != string(user.Role) {
		log.Printf("auth.VerifyAuth: role from token is invalid: token role '%s' user role: '%s' user uid '%s'", tokenInfo.Role, user.Role, user.Uid)
		return tokenInfo, ErrRoleNotMatch
	}

	return tokenInfo, nil
}

var (
	ErrRoleNotMatch        = fmt.Errorf("role not match")
	ErrPermissionsNotMatch = fmt.Errorf("permissions not match")
)
