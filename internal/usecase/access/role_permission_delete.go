package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) DeletePermissionsFromRole(ctx context.Context, role string, permissions []string) error {
	log.Printf("usecasePermissions.DeletePermissionsFromRole: role %s permissions %v", role, permissions)

	if len(role) == 0 {
		log.Printf("usecasePermissions.DeletePermissionsFromRole: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if len(permissions) == 0 {
		log.Printf("usecasePermissions.DeletePermissionsFromRole: empty permissions")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty permissions")
	}

	if err := uc.accessRepository.DeletePermissionsFromRole(ctx, role, permissions); err != nil {
		log.Printf("usecasePermissions.DeletePermissionsFromRole: failed to delete permissions %v from role %s: %v", permissions, role, err)
		return err
	}

	return nil
}
