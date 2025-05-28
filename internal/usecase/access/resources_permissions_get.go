package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetResourcesPermissions(ctx context.Context, limit int64, offset int64) ([]entity.ResourcePermissions, error) {
	log.Printf("usecasePermissions.GetResourcesPermissions: limit %d, offset %d", limit, offset)

	permissions, err := uc.accessRepository.GetResourcesPermissions(ctx, limit, offset)
	if err != nil {
		log.Printf("usecasePermissions.GetResourcesPermissions: failed to get resources permissions: %v", err)
		return nil, err
	}

	if len(permissions) == 0 {
		log.Printf("usecasePermissions.GetResourcesPermissions: no resources permissions found")
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "resources permissions")
	}

	return permissions, nil
}
