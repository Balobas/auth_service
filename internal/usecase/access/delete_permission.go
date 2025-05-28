package useCaseAccess

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (uc *UseCaseAccess) DeletePermission(ctx context.Context, key string) error {
	log.Printf("usecasePermissions.DeletePermission: key %s", key)

	if err := uc.accessRepository.DeletePermission(ctx, key); err != nil {
		log.Printf("usecasePermissions.DeletePermission: failed to delete permission %s: %v", key, err)
		return errors.WithStack(err)
	}

	return nil
}
