package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAccess) GetUserRoles(ctx context.Context, userUid uuid.UUID) ([]entity.Role, error) {
	log.Printf("usecasePemrissions.GetUserRoles: user %s", userUid)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecasePemrissions.GetUserRoles: empty user uid")
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	roles, err := uc.accessRepository.GetUserRoles(ctx, userUid)
	if err != nil {
		log.Printf("usecasePemrissions.GetUserRoles: failed to get user roles: %v", err)
		return nil, err
	}

	if len(roles) == 0 {
		log.Printf("usecasePemrissions.GetUserRoles: no user %s roles found", userUid)
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "user roles")
	}

	return roles, nil
}
