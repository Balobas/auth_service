package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) CreatePermission(ctx context.Context, perm entity.Permission) error {
	log.Printf("usecasePermissions.CreatePermission: key %s, description %s", perm.Key, perm.Description)

	if len(perm.Key) == 0 {
		log.Printf("usecasePermissions.CreatePermission: empty permission key")
		return errors.New("empty permission key")
	}

	if err := uc.permsRepository.CreatePermission(ctx, perm); err != nil {
		log.Printf("usecasePermissions.CreatePermission: failed to create permission %s: %v", perm.Key, err)
		return errors.WithStack(err)
	}
	return nil
}
