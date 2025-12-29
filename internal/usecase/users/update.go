package useCaseUsers

import (
	"context"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/balobas/auth_service/pkg/validations"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
)

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

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		if needUpdateEmail {
			oldUser.Email = user.Email
			oldUser.IsVerified = false
			oldUser.UpdatedAt = time.Now()

			if err := uc.usersRepo.UpdateUser(ctx, oldUser); err != nil {
				log.Printf("usecaseUsers.UpdateUser: failed to update user %s: %v", user.Uid, err)
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
