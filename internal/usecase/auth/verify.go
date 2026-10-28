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

	// TODO: switch tokenInfo.Type: verifyUserToken/verifySystemToken
	switch tokenInfo.Type {
	case entity.TokenTypeUser:
		if err := uc.verifyUserToken(ctx, tokenInfo); err != nil {
			return entity.TokenInfo{}, err
		}
		// Отметка активности для DAU/MAU (неблокирующе, best-effort).
		if uc.activity != nil {
			uc.activity.Touch(tokenInfo.UserUid)
		}
		return tokenInfo, nil
	default:
		return entity.TokenInfo{}, errors.New("unknow token type")
	}
}

func (uc *UseCaseAuth) verifyUserToken(ctx context.Context, tokenInfo entity.TokenInfo) error {
	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("usecaseAuth.VerifyAuth: token expired for user %s", tokenInfo.UserUid)
		return serviceErrors.ErrTokenExpired
	}

	// По uid, а не email: у Telegram-аккаунтов email пустой, GetUserByEmail
	// на нём падает с "empty email". Uid в токене есть всегда.
	user, err := uc.ucUsers.GetUserByUid(ctx, tokenInfo.UserUid)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrNotFound) {
			log.Printf("usecaseAuth.VerifyAuth: user with uid %s not found", tokenInfo.UserUid)
			return errors.Wrap(serviceErrors.ErrInvalidToken, "user not found")
		}
		return err
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		log.Printf("usecaseAuth.VerifyAuth: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "user in token is not user in request")
	}

	roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get user %s roles: %v", user.Uid, err)
		return errors.WithStack(err)
	}
	user.Roles = entity.RolesToStrings(roles)

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAuth: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseAuth.VerifyAuth: session %s not found", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "session not found")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("usecaseAuth.VerifyAuth: token from request already invalid for session %s", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "you should use last issued token")
	}

	if !uuid.Equal(tokenInfo.DeviceUid, session.DeviceUid) {
		log.Printf("usecaseAuth.VerifyAuth: invalid device uid for session %s", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "invalid device in token")
	}

	if len(tokenInfo.Roles) != len(user.Roles) {
		log.Printf("usecaseAuth.VerifyAuth: invalid roles in user %s token: len mismatch", tokenInfo.UserUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "invalid roles in token")
	}

	rolesMap := make(map[string]struct{}, len(tokenInfo.Roles))
	for _, role := range tokenInfo.Roles {
		rolesMap[role] = struct{}{}
	}

	for _, role := range user.Roles {
		if _, ok := rolesMap[role]; !ok {
			log.Printf("usecaseAuth.VerifyAuth: invalid roles in user %s token: role %s doesnt exists in user", tokenInfo.UserUid, role)
			return errors.Wrap(serviceErrors.ErrInvalidToken, "invalid roles in token")
		}
	}

	return nil
}
