package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetRoles(ctx context.Context, rolePattern string, limit int64, offset int64) ([]entity.Role, error) {
	log.Printf("usecasePemrissions.GetRoles: rolePattern %s, limit %d, offset %d", rolePattern, limit, offset)

	roles, err := uc.accessRepository.GetRoles(ctx, rolePattern, limit, offset)
	if err != nil {
		log.Printf("usecasePemrissions.GetRoles: failed to get roles: %v", err)
		return nil, err
	}

	if len(roles) == 0 {
		log.Printf("usecasePemrissions.GetRoles: no roles found")
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "roles")
	}

	return roles, nil
}
