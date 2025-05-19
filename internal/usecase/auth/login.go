package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) Login(ctx context.Context, params entity.LoginParams) (string, string, error) {
	log.Printf("usecaseAuth.Login: email %s", params.Email)

	if len(params.Email) == 0 {
		log.Printf("usecaseAuth.Login: empty email")
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty email"))
	}
	if len(params.Password) == 0 {
		log.Printf("usecaseAuth.Login: empty password")
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty password"))
	}

	if err := validations.ValidateEmail(params.Email); err != nil {
		log.Printf("usecaseAuth.Login: invalid email: %v", err)
		return emptyTokensWithError(errors.Wrapf(serviceErrors.ErrBadRequest, "invalid email: %v", err.Error()))
	}

	user, err := uc.ucUsers.GetUserByEmail(ctx, params.Email)
	if err != nil {
		log.Printf("usecaseAuth.Login: failed to get user by email %s: %v", params.Email, err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	permissions, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseAuth.Login: failed to get user %s permissions: %v", user.Uid, err)
		return emptyTokensWithError(errors.WithStack(err))
	}
	user.Permissions = permissions

	if err := uc.ucCredentials.Validate(ctx, user.Uid, params.Password); err != nil {
		log.Printf("usecaseAuth.Login: wrong password %v", err)
		return emptyTokensWithError(errors.Wrap(err, "wrong password"))
	}

	// _, isFound, err = uc.sessionsRepo.GetSessionByUserUid(ctx, user.Uid)
	// if err != nil {
	// 	log.Printf("failed to get session %v", err)
	// 	return emptyTokensWithError(errors.WithStack(err))
	// }
	// TODO: Return when many sessions are ready
	// if isFound {
	// 	log.Printf("user already authorized")
	// 	return emptyTokensWithError(errors.New("user already authorized"))
	// }

	loginTime := time.Now()

	session := entity.Session{
		Uid:            uuid.NewV4(),
		UserUid:        user.Uid,
		CreatedAt:      loginTime,
		TokensIssuedAt: loginTime.Unix(),
	}

	tokenInfo := entity.TokenInfo{
		UserUid:     user.Uid,
		Email:       user.Email,
		Permissions: user.PermissionsStrings(),
		Role:        string(user.Role),
		SessionUid:  session.Uid,
		IssuedAt:    loginTime.Unix(),
	}

	access, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		log.Printf("usecaseAuth.Login: failed to build jwt for user %s", user.Uid)
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}
	refresh, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		log.Printf("usecaseAuth.Login: failed to build jwt for user %s", user.Uid)
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}

	if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
		log.Printf("usecaseAuth.Login: failed to create session for user %s: %v", user.Uid, err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}

func emptyTokensWithError(err error) (string, string, error) {
	return "", "", err
}
