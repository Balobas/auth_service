package useCaseAuth

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseAuth) VerifyAccess(ctx context.Context, uri string, method string, token string) error {
	log.Printf("usecaseAuth.VerifyAccess: uri: %s, method: %s, token: %s", uri, method, token)

	tokenInfo, err := uc.VerifyAuth(ctx, token)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAccess: failed to verify auth: %v", err)
		return err
	}

	if entity.HasRole(tokenInfo.Roles, entity.UserRoleAdmin) {
		log.Printf("usecaseAuth.VerifyAccess: user is admin")
		return nil
	}

	hasPerms, err := uc.accessRepo.IsUserHasPermissionsForResource(ctx, tokenInfo.UserUid, uri, method)
	if err != nil {
		log.Printf("usecaseAuth.VerifyAccess: failed to check if user has permissions for resource: %v", err)
		return err
	}

	if !hasPerms {
		log.Printf("usecaseAuth.VerifyAccess: user %s has no permissions for resource %s %s", tokenInfo.UserUid, uri, method)
		return errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "user has no permissions for resource")
	}

	return nil
}
