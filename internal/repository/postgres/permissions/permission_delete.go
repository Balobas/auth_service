package repositoryPermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *PermissionsRepository) DeletePermission(ctx context.Context, key string) error {

	row := pgEntity.NewPermissionRow().FromEntity(entity.Permission{Key: key})

	if err := r.Delete(ctx, row, row.ConditionUidEqual()); err != nil {
		log.Printf("repositoryPermissions.DeletePermission: failed to delete permission %s: %v", key, err)
		return errors.WithStack(err)
	}
	return nil
}
