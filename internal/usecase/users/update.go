package useCaseUsers

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
)

//TODO: переосмыслить, мб запихнуть все в одну транзакцию

func (uc *UseCaseUsers) UpdateUser(ctx context.Context, user entity.User, password string) error {
	log.Printf("usecaseUsers.UpdateUser: user uid %s", user.Uid)

	oldUser, isFound, err := uc.usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		log.Printf("usecaseUsers.UpdateUser: failed to get user %s: %v", user.Uid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		return errors.Wrap(serviceErrors.ErrNotFound, "user")
	}

	if len(user.Email) != 0 && oldUser.Email != user.Email {
		if err := validations.ValidateEmail(user.Email); err != nil {
			return errors.Wrap(serviceErrors.ErrBadRequest, err.Error())
		}

		if _, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email); err == nil {
			if isFound {
				return errors.Wrap(serviceErrors.ErrAlreadyExists, "user with email already exists")
			}
		} else {
			log.Printf("usecaseUsers.UpdateUser: failed to get user by email %s: %v", user.Email, err)
			return errors.WithStack(err)
		}
	}

	needUpdateEmail := len(user.Email) != 0 && oldUser.Email != user.Email
	needUpdatePassword := len(password) != 0 && uc.ucCredentials.Validate(ctx, user.Uid, password) != nil

	if !needUpdateEmail && !needUpdatePassword {
		return errors.Wrap(serviceErrors.ErrIdempotentOperation, "nothing to update")
	}

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {
		if needUpdateEmail {
			oldUser.Email = user.Email
			// TODO: подумать над перезаписью пермишенов, так как у юзера слетают все пермишены что были
			oldUser.Permissions = []entity.UserPermission{entity.UserPermissionNotVerified}
			oldUser.UpdatedAt = time.Now()

			if err := uc.usersRepo.UpdateUser(ctx, oldUser); err != nil {
				log.Printf("usecaseUsers.UpdateUser: failed to update user %s: %v", user.Uid, err)
				return err
			}

			if err := uc.permsRepo.UpdateUserPermissions(ctx, oldUser.Uid, oldUser.Permissions); err != nil {
				log.Printf("usecaseUsers.UpdateUser: failed to update user %s permissions: %v", user.Uid, err)
				return err
			}

			if err := uc.ucVerification.CreateVerification(ctx, oldUser.Uid, user.Email); err != nil {
				log.Printf("usecaseUsers.UpdateUser: failed to create verification for user %s: %v", user.Uid, err)
				return err
			}
		}
		if needUpdatePassword {
			if err := uc.ucCredentials.Update(ctx, oldUser.Uid, password); err != nil {
				log.Printf("usecaseUsers.UpdateUser: failed to update user %s password: %v", user.Uid, err)
				return err
			}
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
