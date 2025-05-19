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

func (uc *UseCaseAuth) UpdateUserCreds(ctx context.Context, user entity.User, password string) (string, string, error) {
	log.Printf("usecaseAuth.UpdateUserCreds: user %s", user.Uid)

	if uuid.Equal(user.Uid, uuid.UUID{}) {
		log.Printf("usecaseAuth.UpdateUserCreds: empty user uid")
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid"))
	}
	if len(user.Email) == 0 && len(password) == 0 {
		log.Printf("usecaseAuth.UpdateUserCreds: empty user %s email and password", user.Uid)
		return emptyTokensWithError(errors.Wrap(serviceErrors.ErrBadRequest, "empty user email and password"))
	}

	var access, refresh string

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {

		if err := uc.ucUsers.UpdateUser(ctx, user, password); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to update user %s: %v", user.Uid, err)
			return err
		}

		perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
		if err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to get user %s permissions: %v", user.Uid, err)
			return err
		}
		user.Permissions = perms

		now := time.Now()

		session := entity.Session{
			Uid:            uuid.NewV4(),
			UserUid:        user.Uid,
			TokensIssuedAt: now.Unix(),
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		tokenInfo := entity.TokenInfo{
			UserUid:     user.Uid,
			Email:       user.Email,
			Permissions: user.PermissionsStrings(),
			Role:        string(user.Role),
			SessionUid:  session.Uid,
			IssuedAt:    now.Unix(),
		}

		if err := uc.sessionsRepo.DeleteSessionByUserUid(ctx, user.Uid); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to delete old session for user %s: %v", user.Uid, err)
			return err
		}

		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to create new session for user %s: %v", user.Uid, err)
			return err
		}

		access, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
		if err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to build jwt for user %s: %v", user.Uid, err)
			return errors.Wrapf(err, "failed to build jwt")
		}
		refresh, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
		if err != nil {
			log.Printf("usecaseAuth.UpdateUserCreds: failed to build jwt for user %s: %v", user.Uid, err)
			return errors.Wrapf(err, "failed to build jwt")
		}

		return nil
	}); err != nil {
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
