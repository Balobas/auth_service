package useCaseAuth

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) VerifyAuth(ctx context.Context, token string) (entity.TokenInfo, error) {
	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		return entity.TokenInfo{}, errors.New("token expired")
	}

	user, isFound, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		return entity.TokenInfo{}, errors.New("user not found")
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	if tokenInfo.Role != string(user.Role) {
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	if len(tokenInfo.Permissions) != len(user.Permissions) {
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	permsMap := make(map[entity.UserPermission]struct{}, len(tokenInfo.Permissions))
	for _, perm := range tokenInfo.Permissions {
		permsMap[entity.UserPermission(perm)] = struct{}{}
	}

	for _, perm := range user.Permissions {
		if _, ok := permsMap[perm]; !ok {
			return entity.TokenInfo{}, errors.New("invalid token")
		}
	}

	_, isFound, err = uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if !isFound {
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	return tokenInfo, nil
}
