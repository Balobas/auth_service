package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) UpdatePermission(ctx context.Context, perm entity.Permission) error {
	log.Printf("usecasePermissions.UpdatePermission: key %s description %s", perm.Key, perm.Description)

	if len(perm.Key) == 0 {
		log.Printf("usecasePermissions.UpdatePermission: empty permission key")
		return errors.New("empty permission key")
	}

	if err := uc.permsRepository.UpdatePermission(ctx, perm); err != nil {
		log.Printf("usecasePermissions.UpdatePermission: failed to update permission %s: %v", perm.Key, err)
		return errors.WithStack(err)
	}
	return nil
}
