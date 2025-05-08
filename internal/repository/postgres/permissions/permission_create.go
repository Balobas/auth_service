package repositoryPermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *PermissionsRepository) CreatePermission(ctx context.Context, perm entity.Permission) error {

	if err := r.Create(ctx, pgEntity.NewPermissionRow().FromEntity(perm)); err != nil {
		log.Printf("repositoryPermissions.CreatePermission: failed to create permission (key %s): %v", perm.Key, err)
		return errors.WithStack(err)
	}
	return nil
}
