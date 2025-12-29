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

	user.Uid = uuid.NewV4()
	user.IsVerified = false
	user.CreatedAt = time.Now()

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {

		_, isFound, err := uc.usersRepo.GetByEmail(ctx, user.Email)
		if err != nil {
			log.Printf("usecaseUsers.Register: failed to get user by email %s: %v", user.Email, err)
			return err
		}
		if isFound {
			log.Printf("usecaseUsers.Register: user with email %s already exist", user.Email)
			return errors.Wrap(serviceErrors.ErrAlreadyExists, "user with email already exists")
		}

		if err := uc.usersRepo.CreateUser(ctx, user); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user (email %s): %v", user.Email, err)
			return err
		}

		if err := uc.accessRepo.AddRoleToUser(ctx, user.Uid, entity.UserRoleUser); err != nil {
			log.Printf("usecaseUsers.Register: failed to add role user to user (email %s): %v", user.Email, err)
			return err
		}

		if err := uc.ucCredentials.Create(ctx, user.Uid, password); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user (email %s) credentials: %v", user.Email, err)
			return err
		}

		if err := uc.ucVerification.CreateVerification(ctx, user.Uid, user.Email); err != nil {
			log.Printf("usecaseUsers.Register: failed to create verification for user (email %s): %v", user.Email, err)
			return err
		}

		if err := uc.ucOutboxMessages.CreateUserRegisteredMessage(ctx, user); err != nil {
			log.Printf("usecaseUsers.Register: failed to create user registered message (email %s): %v", user.Email, err)
			return err
		}

		return nil
	}); err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return user.Uid, nil
}
