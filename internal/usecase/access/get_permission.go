package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetPermission(ctx context.Context, key string) (entity.Permission, error) {
	log.Printf("usecasePermissions.GetPermission: key %s", key)

	perm, isFound, err := uc.accessRepository.GetPermission(ctx, key)
	if err != nil {
		log.Printf("usecasePermissions.GetPermission: failed to get permission %s: %v", key, err)
		return entity.Permission{}, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecasePermissions.GetPermission: permission %s not found", key)
		return entity.Permission{}, errors.Wrap(serviceErrors.ErrNotFound, "permission")
	}

	return perm, nil
}
