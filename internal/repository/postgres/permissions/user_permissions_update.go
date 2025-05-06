package repositoryPermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid, Permissions: perms})
	if err := r.Update(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to update permissions %v", err)
		return errors.Wrapf(err, "failed to update permissions for user %s", userUid)
	}

	log.Printf("successfully update permissions")
	return nil
}
