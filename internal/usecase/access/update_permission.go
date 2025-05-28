package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) UpdatePermission(ctx context.Context, perm entity.Permission) error {
	log.Printf("usecasePermissions.UpdatePermission: key %s description %s", perm.Key, perm.Description)

	if len(perm.Key) == 0 {
		log.Printf("usecasePermissions.UpdatePermission: empty permission key")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty permission key")
	}

	if err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		_, isFound, err := uc.accessRepository.GetPermission(ctx, perm.Key)
		if err != nil {
			log.Printf("usecasePermissions.UpdatePermission: failed to get permission %s: %v", perm.Key, err)
			return err
		}
		if !isFound {
			log.Printf("usecasePermissions.UpdatePermission: permission %s not found", perm.Key)
			return errors.Wrap(serviceErrors.ErrNotFound, "permission")
		}

		if err := uc.accessRepository.UpdatePermission(ctx, perm); err != nil {
			log.Printf("usecasePermissions.UpdatePermission: failed to update permission %s: %v", perm.Key, err)
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
