package useCasePermissions

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) DeletePermission(ctx context.Context, key string) error {
	log.Printf("usecasePermissions.DeletePermission: key %s", key)

	err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		if err := uc.permsRepository.DeletePermission(ctx, key); err != nil {
			log.Printf("usecasePermissions.DeletePermission: failed to delete permission %s: %v", key, err)
			return err
		}

		if err := uc.permsRepository.AddPermissionToDeletingList(ctx, key); err != nil {
			log.Printf("usecasePermissions.DeletePermission: failed to add permission %s to deleting list: %v", key, err)
			return err
		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
