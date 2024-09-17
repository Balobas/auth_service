package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) Login(ctx context.Context, params entity.LoginParams) (string, string, error) {
	log.Printf("auth.Login email: %s\n", params.Email)

	if err := validations.ValidateEmail(params.Email); err != nil {
		log.Printf("invalid email: %v", err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	if len(params.Password) < uc.cfg.MinPasswordLen() {
		log.Printf("password shoud have >= %d symbols", uc.cfg.MinPasswordLen())
		return emptyTokensWithError(errors.Errorf("password shoud have >= %d symbols", uc.cfg.MinPasswordLen()))
	}

	user, isFound, err := uc.ucUsers.GetUserByEmail(ctx, params.Email)
	if err != nil {
		log.Printf("failed to get user by email")
		return emptyTokensWithError(errors.WithStack(err))
	}
	if !isFound {
		log.Printf("user not found")
		return emptyTokensWithError(errors.New("user not found"))
	}

	permissions, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
	if err != nil {
		log.Printf("failed to get user permissions")
		return emptyTokensWithError(errors.WithStack(err))
	}
	user.Permissions = permissions

	if err := uc.ucCredentials.Validate(ctx, user.Uid, params.Password); err != nil {
		log.Printf("wrong password %v", err)
		return emptyTokensWithError(errors.New("wrong password"))
	}

	_, isFound, err = uc.sessionsRepo.GetSessionByUserUid(ctx, user.Uid)
	if err != nil {
		log.Printf("failed to get session %v", err)
		return emptyTokensWithError(errors.WithStack(err))
	}
	if isFound {
		log.Printf("user already authorized")
		return emptyTokensWithError(errors.New("user already authorized"))
	}

	session := entity.Session{
		Uid:       uuid.NewV4(),
		UserUid:   user.Uid,
		CreatedAt: time.Now(),
	}

	tokenInfo := entity.TokenInfo{
		UserUid:     user.Uid,
		Email:       user.Email,
		Permissions: user.PermissionsStrings(),
		Role:        string(user.Role),
		SessionUid:  session.Uid,
	}

	access, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
	if err != nil {
		log.Printf("failed to build jwt")
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}
	refresh, err := uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
	if err != nil {
		log.Printf("failed to build jwt")
		return emptyTokensWithError(errors.Wrapf(err, "failed to build jwt"))
	}

	if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
		log.Printf("failed to create session %v", err)
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}

func emptyTokensWithError(err error) (string, string, error) {
	return "", "", err
}
