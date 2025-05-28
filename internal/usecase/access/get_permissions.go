package useCaseAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error) {
	log.Printf("usecasePermissions.GetPermissions: key pattern %s", keyPattern)

	perms, err := uc.accessRepository.GetPermissions(ctx, keyPattern)
	if err != nil {
		log.Printf("usecasePermissions.GetPermissions: failed to get permission (key pattern %s): %v", keyPattern, err)
		return nil, errors.WithStack(err)
	}

	if len(perms) == 0 {
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "permissions")
	}

	return perms, nil
}
