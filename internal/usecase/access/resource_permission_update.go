package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) UpdateResourcePermission(ctx context.Context, uri string, method string, permission string) error {
	log.Printf("usecasePermissions.UpdateResourcePermission: uri %s, method %s, permission %s", uri, method, permission)

	if len(uri) == 0 {
		log.Printf("usecasePermissions.UpdateResourcePermission: empty uri")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty uri")
	}

	if len(method) == 0 {
		log.Printf("usecasePermissions.UpdateResourcePermission: empty method")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty method")
	}

	if len(permission) == 0 {
		log.Printf("usecasePermissions.UpdateResourcePermission: empty permission")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty permission")
	}

	// TODO: потестить
	if err := uc.accessRepository.UpdateResourcePermission(ctx, uri, method, permission); err != nil {
		log.Printf("usecasePermissions.UpdateResourcePermission: failed to update permission %s for resource %s method %s: %v", permission, uri, method, err)
		return err
	}

	return nil
}
