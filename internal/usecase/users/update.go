package useCaseUsers

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
)

func (uc *UseCaseUsers) UpdateUser(ctx context.Context, user entity.User, password string) error {
	if len(user.Email) != 0 {
		if err := validations.ValidateEmail(user.Email); err != nil {
			return errors.WithStack(err)
		}
	}

	oldUser, isFound, err := uc.usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isFound {
		return errors.New("user not found")
	}

	needUpdateEmail := len(user.Email) != 0 && oldUser.Email != user.Email
	needUpdatePassword := len(password) != 0 && uc.ucCredentials.Validate(ctx, user.Uid, password) != nil

	if !needUpdateEmail && !needUpdatePassword {
		return errors.New("nothing to update")
	}

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {
		if needUpdateEmail {
			oldUser.Email = user.Email
			oldUser.Permissions = []entity.UserPermission{entity.UserPermissionNotVerified}
			oldUser.UpdatedAt = time.Now()

			if err := uc.usersRepo.UpdateUser(ctx, oldUser); err != nil {
				return err
			}

			if err := uc.permsRepo.UpdateUserPermissions(ctx, oldUser.Uid, oldUser.Permissions); err != nil {
				return err
			}

			if err := uc.ucVerification.CreateVerification(ctx, oldUser.Uid, user.Email); err != nil {
				return err
			}
		}
		if needUpdatePassword {
			if err := uc.ucCredentials.Update(ctx, oldUser.Uid, password); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
