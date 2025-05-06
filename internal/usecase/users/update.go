package useCaseUsers

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
)

func (uc *UseCaseUsers) UpdateUser(ctx context.Context, user entity.User, password string) error {
	oldUser, isFound, err := uc.usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isFound {
		return errors.New("user not found")
	}

	if len(user.Email) != 0 && oldUser.Email != user.Email {
		if err := validations.ValidateEmail(user.Email); err != nil {
			return errors.WithStack(err)
		}

		if _, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email); err == nil {
			if isFound {
				return errors.New("user with email already exists")
			}
		} else {
			return errors.WithStack(err)
		}
	}

	needUpdateEmail := len(user.Email) != 0 && oldUser.Email != user.Email
	needUpdatePassword := len(password) != 0 && uc.ucCredentials.Validate(ctx, user.Uid, password) != nil

	if !needUpdateEmail && !needUpdatePassword {
		return nil
	}

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {
		if needUpdateEmail {
			oldUser.Email = user.Email
			// TODO: подумать над перезаписью пермишенов, так как у юзера слетают все пермишены что были
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
