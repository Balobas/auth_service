package useCaseAuth

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) UpdateUserCreds(ctx context.Context, user entity.User, password string) (string, string, error) {
	log.Printf("auth.UpdateUserCreds")

	var access, refresh string

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {

		if err := uc.ucUsers.UpdateUser(ctx, user, password); err != nil {
			log.Printf("failed to update user: %v", errors.WithStack(err))
			return err
		}

		perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
		if err != nil {
			log.Printf("failed to get permissions: %v", err)
			return err
		}
		user.Permissions = perms

		session := entity.Session{
			Uid:       uuid.NewV4(),
			UserUid:   user.Uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tokenInfo := entity.TokenInfo{
			UserUid:     user.Uid,
			Email:       user.Email,
			Permissions: user.PermissionsStrings(),
			Role:        string(user.Role),
			SessionUid:  session.Uid,
		}

		if err := uc.sessionsRepo.DeleteSessionByUserUid(ctx, user.Uid); err != nil {
			log.Printf("failed to delete old session: %v", err)
			return err
		}

		if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
			log.Printf("failed to create new session: %v", err)
			return err
		}

		access, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.AccessJwtTTL())
		if err != nil {
			log.Printf("failed to build jwt")
			return errors.Wrapf(err, "failed to build jwt")
		}
		refresh, err = uc.jwtManager.NewToken(tokenInfo, uc.cfg.RefreshJwtTTL())
		if err != nil {
			log.Printf("failed to build jwt")
			return errors.Wrapf(err, "failed to build jwt")
		}

		return nil
	}); err != nil {
		log.Printf("failed to update user creds")
		return emptyTokensWithError(errors.WithStack(err))
	}

	return access, refresh, nil
}
