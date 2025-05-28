package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) AddPermissionsToRole(ctx context.Context, role string, permissions []string) error {
	log.Printf("usecasePemrissions.AddPermissionsToRole: role %s permissions %v", role, permissions)

	if len(role) == 0 {
		log.Printf("usecasePemrissions.AddPermissionsToRole: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if len(permissions) == 0 {
		log.Printf("usecasePemrissions.AddPermissionsToRole: empty permissions")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty permissions")
	}

	// Роль и пермишены проверятся в БД, там constraint
	if err := uc.accessRepository.AddPermissionsToRole(ctx, role, permissions); err != nil {
		log.Printf("usecasePemrissions.AddPermissionsToRole: failed to add permissions %v to role %s: %v", permissions, role, err)
		return err
	}

	return nil
}
