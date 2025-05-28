package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetRolePermissions(ctx context.Context, role string) ([]entity.Permission, error) {
	log.Printf("usecasePermissions.GetRolePermissions: role %s", role)

	if len(role) == 0 {
		log.Printf("usecasePermissions.GetRolePermissions: empty role")
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	permissions, err := uc.accessRepository.GetRolePermissions(ctx, role)
	if err != nil {
		log.Printf("usecasePermissions.GetRolePermissions: failed to get role %s permissions: %v", role, err)
		return nil, err
	}

	if len(permissions) == 0 {
		log.Printf("usecasePermissions.GetRolePermissions: no permissions found for role %s", role)
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "permissions")
	}

	return permissions, nil
}
