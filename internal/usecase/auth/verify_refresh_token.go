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

func (uc *UseCaseAuth) parseAndVerifyRefreshToken(ctx context.Context, token string) (entity.TokenInfo, error) {
	log.Printf("usecaseAuth.parseAndVerifyRefreshToken, token: %v", token)

	if len(token) == 0 {
		log.Printf("usecaseAuth.parseAndVerifyRefreshToken: empty token")
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrBadRequest, "empty refresh token")
	}

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("usecaseAuth.parseAndVerifyRefreshToken: failed to parse token\n")
		return entity.TokenInfo{}, errors.Wrap(serviceErrors.ErrBadRequest, err.Error())
	}
	// TODO: switch tokenInfo.Type: verifyUserToken/verifySystemToken
	switch tokenInfo.Type {
	case entity.TokenTypeUser:
		return tokenInfo, uc.verifyUserRefreshToken(ctx, tokenInfo)
	case entity.TokenTypeSystem:
		return tokenInfo, uc.verifySystemRefreshToken(ctx, tokenInfo)
	default:
		return entity.TokenInfo{}, errors.New("unknow token type")
	}
}

func (uc *UseCaseAuth) verifyUserRefreshToken(ctx context.Context, tokenInfo entity.TokenInfo) error {
	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("usecaseAuth.verifyUserRefreshToken: token expired for user %s", tokenInfo.UserUid)
		return serviceErrors.ErrTokenExpired
	}

	user, err := uc.ucUsers.GetUserByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("usecaseAuth.verifyUserRefreshToken: failed to get user %s by email %s: %v\n", tokenInfo.UserUid, tokenInfo.Email, err)
		return errors.WithStack(err)
	}

	if !uuid.Equal(user.Uid, tokenInfo.UserUid) {
		// на случаи когда в токене email реального юзера, а uid не реального
		log.Printf("usecaseAuth.verifyUserRefreshToken: user in token (%s) is not user in request (%s)", tokenInfo.UserUid, user.Uid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "user in token is not user in request")
	}

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("usecaseAuth.verifyUserRefreshToken: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseAuth.verifyUserRefreshToken: session %s not found", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrNotFound, "session")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("usecaseAuth.verifyUserRefreshToken: token from request already invalid for session %s", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "you should use last issued token")
	}
	if !uuid.Equal(session.DeviceUid, tokenInfo.DeviceUid) {
		return errors.Wrap(serviceErrors.ErrInvalidToken, "wrong device")
	}

	return nil
}

func (uc *UseCaseAuth) verifySystemRefreshToken(ctx context.Context, tokenInfo entity.TokenInfo) error {
	if tokenInfo.ExpiredAt <= time.Now().Unix() {
		log.Printf("usecaseAuth.verifySystemRefreshToken: token expired for service %s", tokenInfo.ServiceUid)
		return serviceErrors.ErrTokenExpired
	}

	service, isFound, err := uc.servicesRepo.GetServiceByUid(ctx, tokenInfo.ServiceUid)
	if err != nil {
		log.Printf("usecaseAuth.verifySystemRefreshToken: failed to get service %s by  uid: %v\n", tokenInfo.ServiceUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		return errors.Wrap(serviceErrors.ErrNotFound, "service")
	}

	if service.Name != tokenInfo.ServiceName {
		log.Printf("usecaseAuth.verifySystemRefreshToken: service name in token (%s) is not service name in request (%s)", tokenInfo.ServiceName, service.Name)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "invalid service name")
	}
	if service.Domain != tokenInfo.Domain {
		log.Printf("usecaseAuth.verifySystemRefreshToken: service domain in token (%s) is not service domain in request (%s)", tokenInfo.Domain, service.Domain)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "invalid service domain")
	}

	session, isFound, err := uc.sessionsRepo.GetSessionByUid(ctx, tokenInfo.SessionUid)
	if err != nil {
		log.Printf("usecaseAuth.verifySystemRefreshToken: failed to get session %s: %v", tokenInfo.SessionUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseAuth.verifySystemRefreshToken: session %s not found", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrNotFound, "session")
	}

	if session.TokensIssuedAt != tokenInfo.IssuedAt {
		log.Printf("usecaseAuth.verifySystemRefreshToken: token from request already invalid for session %s", tokenInfo.SessionUid)
		return errors.Wrap(serviceErrors.ErrInvalidToken, "you should use last issued token")
	}
	if !uuid.Equal(session.DeviceUid, tokenInfo.DeviceUid) {
		return errors.Wrap(serviceErrors.ErrInvalidToken, "wrong device")
	}

	return nil
}
