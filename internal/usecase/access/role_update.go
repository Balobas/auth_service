package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) UpdateRole(ctx context.Context, role entity.Role) error {
	log.Printf("usecasePermissions.UpdateRole: role %s", role.Role)

	if len(role.Role) == 0 {
		log.Printf("usecasePermissions.UpdateRole: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if err := uc.accessRepository.UpdateRole(ctx, role); err != nil {
		log.Printf("usecasePermissions.UpdateRole: failed to update role %s: %v", role.Role, err)
		return errors.WithStack(err)
	}
	return nil
}
