package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetRole(ctx context.Context, role string) (entity.Role, error) {
	log.Printf("usecasePemrissions.GetRole: role %s", role)

	if len(role) == 0 {
		log.Printf("usecasePemrissions.GetRole: empty role")
		return entity.Role{}, errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	res, isFound, err := uc.accessRepository.GetRole(ctx, role)
	if err != nil {
		log.Printf("usecasePemrissions.GetRole: failed to get role %s: %v", role, err)
		return entity.Role{}, err
	}

	if !isFound {
		log.Printf("usecasePemrissions.GetRole: role %s not found", role)
		return entity.Role{}, errors.Wrap(serviceErrors.ErrNotFound, "role")
	}

	return res, nil
}
