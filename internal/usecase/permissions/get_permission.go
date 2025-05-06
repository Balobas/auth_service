package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) GetPermission(ctx context.Context, key string) (entity.Permission, bool, error) {
	log.Printf("usecasePermissions.GetPermission: key %s", key)

	perm, isFound, err := uc.permsRepository.GetPermission(ctx, key)
	if err != nil {
		log.Printf("usecasePermissions.GetPermission: failed to get permission %s: %v", key, err)
		return entity.Permission{}, false, errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecasePermissions.GetPermission: permission %s not found", key)
		return entity.Permission{}, false, nil
	}

	return perm, true, nil
}
