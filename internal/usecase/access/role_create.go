package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) CreateRole(ctx context.Context, role entity.Role) error {
	log.Printf("usecasePermissions.CreateRole: role %s", role.Role)

	if len(role.Role) == 0 {
		log.Printf("usecasePermissions.CreateRole: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if err := uc.accessRepository.CreateRole(ctx, role); err != nil {
		log.Printf("usecasePermissions.CreateRole: failed to create role %s: %v", role.Role, err)
		return errors.WithStack(err)
	}
	return nil
}
