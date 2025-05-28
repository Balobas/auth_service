package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) DeleteRole(ctx context.Context, role string) error {
	log.Printf("usecasePermissions.DeleteRole: role %s", role)

	if len(role) == 0 {
		log.Printf("usecasePermissions.DeleteRole: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if err := uc.accessRepository.DeleteRole(ctx, role); err != nil {
		log.Printf("usecasePermissions.DeleteRole: failed to delete role %s: %v", role, err)
		return errors.WithStack(err)
	}
	return nil
}
