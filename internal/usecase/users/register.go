package useCaseUsers

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error) {

	_, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}
	if isFound {
		return uuid.UUID{}, errors.New("user with email is already exists")
	}

	user.Uid = uuid.NewV4()
	user.Permissions = []entity.UserPermission{entity.UserPermissionNotVerified}
	user.Role = entity.UserRoleUser
	user.CreatedAt = time.Now()

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {

		if err := uc.usersRepo.CreateUser(ctx, user); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.permsRepo.CreateUserPermissions(ctx, user.Uid, user.Permissions); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.ucCredentials.Create(ctx, user.Uid, password); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.ucVerification.CreateVerification(ctx, user.Uid); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return user.Uid, nil
}
