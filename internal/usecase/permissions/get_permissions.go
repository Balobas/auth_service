package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error) {
	log.Printf("usecasePermissions.GetPermissions: key pattern %s", keyPattern)

	perms, err := uc.permsRepository.GetPermissions(ctx, keyPattern)
	if err != nil {
		log.Printf("usecasePermissions.GetPermissions: failed to get permission (key pattern %s): %v", keyPattern, err)
		return nil, errors.WithStack(err)
	}

	return perms, nil
}
