package useCaseUsers

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

func (uc *UseCaseUsers) Register(ctx context.Context, user entity.User, password string) (uuid.UUID, error) {
	log.Printf("usecaseUsers.Register: email %s", user.Email)

	if err := validations.ValidateEmail(user.Email); err != nil {
		log.Printf("usecaseUsers.Register: invalid email %s: %v", user.Email, err)
		return uuid.UUID{}, errors.Wrap(serviceErrors.ErrBadRequest, err.Error())
	}

	if len(password) < uc.cfg.MinPasswordLen() {
		log.Printf("usecaseUsers.Register: password shoud have >= %d symbols", uc.cfg.MinPasswordLen())
		return uuid.UUID{}, errors.Wrapf(serviceErrors.ErrBadRequest, "password shoud have >= %d symbols", uc.cfg.MinPasswordLen())
	}

	_, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("usecaseUsers.Register: failed to get user by email %s: %v", user.Email, err)
		return uuid.UUID{}, errors.WithStack(err)
	}
	if isFound {
		log.Printf("usecaseUsers.Register: user with email %s already exist", user.Email)
		return uuid.UUID{}, errors.Wrap(serviceErrors.ErrAlreadyExists, "user with email is already exists")
	}

	user.Uid = uuid.NewV4()
	user.Permissions = []entity.UserPermission{entity.UserPermissionNotVerified}
	user.CreatedAt = time.Now()

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {

		if err := uc.usersRepo.CreateUser(ctx, user); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user (email %s): %v", user.Email, err)
			return err
		}

		if err := uc.permsRepo.CreateUserPermissions(ctx, user.Uid, user.Permissions); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user (email %s) permissions: %v", user.Email, err)
			return err
		}

		if err := uc.ucCredentials.Create(ctx, user.Uid, password); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user (email %s) credentials: %v", user.Email, err)
			return err
		}

		if user.Role != entity.UserRoleAdmin {
			if err := uc.ucVerification.CreateVerification(ctx, user.Uid, user.Email); err != nil {
				log.Printf("usecaseUsers.Register: failed to create verification for user (email %s): %v", user.Email, err)
				return err
			}

			if err := uc.ucOutboxMessages.CreateUserRegisteredMessage(ctx, user); err != nil {
				log.Printf("usecaseUsers.Register: failed to create user registered message (email %s): %v", user.Email, err)
				return err
			}
		}

		return nil
	}); err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return user.Uid, nil
}
