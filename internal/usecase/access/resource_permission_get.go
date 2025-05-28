package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) GetResourcePermission(ctx context.Context, uri string, method string) (string, error) {
	log.Printf("usecasePermissions.GetResourcePermission: uri %s, method %s", uri, method)

	if len(uri) == 0 {
		log.Printf("usecasePermissions.GetResourcePermission: empty uri")
		return "", errors.Wrap(serviceErrors.ErrBadRequest, "empty uri")
	}

	if len(method) == 0 {
		log.Printf("usecasePermissions.GetResourcePermission: empty method")
		return "", errors.Wrap(serviceErrors.ErrBadRequest, "empty method")
	}

	permission, isFound, err := uc.accessRepository.GetResourcePermission(ctx, uri, method)
	if err != nil {
		log.Printf("usecasePermissions.GetResourcePermission: failed to get permission for resource %s method %s: %v", uri, method, err)
		return "", err
	}

	if !isFound {
		log.Printf("usecasePermissions.GetResourcePermission: permission not found for resource %s method %s", uri, method)
		return "", errors.Wrap(serviceErrors.ErrNotFound, "permission")
	}

	return permission, nil
}
