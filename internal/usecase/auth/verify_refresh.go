package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) verifyRefreshToken(ctx context.Context, token string) (entity.TokenInfo, error) {
	log.Printf("auth.verifyRefreshToken, token: %v", token)

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("auth.verifyRefreshToken: failed to parse token\n")
		return entity.TokenInfo{}, errors.WithStack(err)
	}

	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("auth.verifyRefreshToken: token expired for user %s", tokenInfo.UserUid)
		return entity.TokenInfo{}, errors.New("token expired")
	}

	user, isFound, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("auth.verifyRefreshToken: failed to get user %s: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("auth.verifyRefreshToken: user %s not found: %v\n", tokenInfo.UserUid, err)
		return entity.TokenInfo{}, errors.New("user not found")
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("auth.verifyRefreshToken: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("auth.verifyRefreshToken: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return entity.TokenInfo{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("auth.verifyRefreshToken: session %s not found", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.New("invalid token")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("auth.verifyRefreshToken: token from request already invalid for session %s", tokenInfo.SessionUid)
		return entity.TokenInfo{}, errors.New("invalid token. you should use last issued token")
	}

	return tokenInfo, nil
}
