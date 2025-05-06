package repositoryPermissions

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *PermissionsRepository) RemovePermissionFromDeletingList(ctx context.Context, key string) error {
	log.Printf("permissionsRepository.RemovePermissionFromDeletingList: permission %s", key)
	if len(key) == 0 {
		log.Printf("permissionsRepository.RemovePermissionFromDeletingList: empty permission")
		return errors.New("empty permission")
	}

	stmt := "delete from deleting_permissions where key=$1"
	if _, err := r.DB().Exec(ctx, stmt, key); err != nil {
		log.Printf("permissionsRepository.RemovePermissionFromDeletingList: failed to remove permission %s: %v", key, err)
		return errors.WithStack(err)
	}
	return nil
}
