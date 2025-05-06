package repositoryPermissions

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *PermissionsRepository) AddPermissionToDeletingList(ctx context.Context, key string) error {
	log.Printf("permissionsRepository.AddPermissionToDeletingList: permission %s", key)
	if len(key) == 0 {
		log.Printf("permissionsRepository.AddPermissionToDeletingList: empty permission")
		return errors.New("empty permission")
	}

	stmt := "insert into deleting_permissions (key) values ($1)"
	if _, err := r.DB().Exec(ctx, stmt, key); err != nil {
		log.Printf("permissionsRepository.AddPermissionToDeletingList: failed to add permission %s: %v", key, err)
		return errors.WithStack(err)
	}
	return nil
}
