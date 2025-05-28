package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) AddResourcePermission(ctx context.Context, uri string, method string, permission string) error {
	log.Printf("usecasePermissions.AddResourcePermission: uri %s, method %s, permission %s", uri, method, permission)

	if len(uri) == 0 {
		log.Printf("usecasePermissions.AddResourcePermission: empty uri")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty uri")
	}

	if len(method) == 0 {
		log.Printf("usecasePermissions.AddResourcePermission: empty method")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty method")
	}

	if len(permission) == 0 {
		log.Printf("usecasePermissions.AddResourcePermission: empty permission")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty permission")
	}

	if err := uc.accessRepository.AddResourcePermission(ctx, uri, method, permission); err != nil {
		log.Printf("usecasePermissions.AddResourcePermission: failed to add permission %s to resource %s method %s: %v", permission, uri, method, err)
		return err
	}

	return nil
}
