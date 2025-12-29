package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAccess) AddRoleToUser(ctx context.Context, userUid uuid.UUID, role string) error {
	log.Printf("usecasePemrissions.AddRoleToUser: user %s role %s", userUid, role)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecasePemrissions.AddRoleToUser: empty user uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	if len(role) == 0 {
		log.Printf("usecasePemrissions.AddRoleToUser: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {

		_, isFound, err := uc.usersRepository.GetUserByUid(ctx, userUid)
		if err != nil {
			log.Printf("usecasePemrissions.AddRoleToUser: failed to get user %s: %v", userUid, err)
			return err
		}

		if !isFound {
			log.Printf("usecasePemrissions.AddRoleToUser: user %s not found", userUid)
			return errors.Wrap(serviceErrors.ErrNotFound, "user")
		}

		role, isFound, err := uc.accessRepository.GetRole(ctx, role)
		if err != nil {
			log.Printf("usecasePemrissions.AddRoleToUser: failed to get role %s: %v", role, err)
			return err
		}

		if !isFound {
			log.Printf("usecasePemrissions.AddRoleToUser: role %s not found", role)
			return errors.Wrap(serviceErrors.ErrNotFound, "role")
		}

		if err := uc.accessRepository.AddRoleToUser(ctx, userUid, role.Role); err != nil {
			log.Printf("usecasePemrissions.AddRoleToUser: failed to add role %s to user %s: %v", role.Role, userUid, err)
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
