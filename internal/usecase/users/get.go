package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, error) {
	log.Printf("usecaseUsers.GetUserByUid: uid %s", uid)

	if uuid.Equal(uid, uuid.UUID{}) {
		log.Printf("usecaseUsers.GetUserByUid: empty user uid")
		return entity.User{}, errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	var (
		user    entity.User
		isFound bool
		err     error
	)
	if err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		user, isFound, err = uc.usersRepo.GetUserByUid(ctx, uid)
		if err != nil {
			log.Printf("usecaseUsers.GetUserByUid: failed to get user %s: %v", uid, err)
			return err
		}
		if !isFound {
			log.Printf("usecaseUsers.GetUserByUid: user %s not found", uid)
			return nil
		}
		isFound = true

		perms, err := uc.permsRepo.GetUserPermissions(ctx, uid)
		if err != nil {
			log.Printf("usecaseUsers.GetUserByUid: failed to get user %s permissions: %v", uid, err)
			return err
		}

		user.Permissions = perms
		return nil
	}); err != nil {
		return entity.User{}, errors.WithStack(err)
	}
	if !isFound {
		return entity.User{}, errors.Wrap(serviceErrors.ErrNotFound, "user")
	}

	return user, nil
}

func (uc *UseCaseUsers) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	log.Printf("usecaseUsers.GetUserByEmail: email %s", email)

	if err := validations.ValidateEmail(email); err != nil {
		return entity.User{}, errors.Wrap(serviceErrors.ErrBadRequest, err.Error())
	}

	var (
		user    entity.User
		isFound bool
		err     error
	)
	if err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		user, isFound, err = uc.usersRepo.GetByEmail(ctx, email)
		if err != nil {
			log.Printf("usecaseUsers.GetUserByEmail: failed to get user with email %s: %v", email, err)
			return err
		}
		if !isFound {
			log.Printf("usecaseUsers.GetUserByEmail: user with email %s not found", email)
			return nil
		}
		isFound = true

		perms, err := uc.permsRepo.GetUserPermissions(ctx, user.Uid)
		if err != nil {
			log.Printf("usecaseUsers.GetUserByEmail: failed to get user %s permissions: %v", user.Uid, err)
			return err
		}

		user.Permissions = perms
		return nil
	}); err != nil {
		return entity.User{}, errors.WithStack(err)
	}
	if !isFound {
		return entity.User{}, errors.Wrap(serviceErrors.ErrNotFound, "user")
	}

	return user, nil
}
