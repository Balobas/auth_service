package useCaseUsers

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error) {
	if err := validations.ValidateEmail(user.Email); err != nil {
		log.Printf("invalid email: %s %v", user.Email, err)
		return uuid.UUID{}, errors.Wrap(err, "invalid email")
	}

	_, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("failed to get user by email: %v\n", errors.WithStack(err))
		return uuid.UUID{}, errors.WithStack(err)
	}
	if isFound {
		log.Printf("user with email %s already exist\n", user.Email)
		return uuid.UUID{}, errors.New("user with email is already exists")
	}

	user.Uid = uuid.NewV4()
	user.Permissions = []entity.UserPermission{entity.UserPermissionNotVerified}
	user.Role = entity.UserRoleUser
	user.CreatedAt = time.Now()

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {

		if err := uc.usersRepo.CreateUser(ctx, user); err != nil {
			log.Printf("failed to create user: %v\n", err)
			return errors.WithStack(err)
		}

		if err := uc.permsRepo.CreateUserPermissions(ctx, user.Uid, user.Permissions); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.ucCredentials.Create(ctx, user.Uid, password); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.ucVerification.CreateVerification(ctx, user.Uid, user.Email); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		log.Printf("failed to create user in tx: %v", errors.WithStack(err))
		return uuid.UUID{}, errors.WithStack(err)
	}

	return user.Uid, nil
}
