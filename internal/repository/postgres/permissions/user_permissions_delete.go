package repositoryPermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) DeleteUserPermissions(ctx context.Context, userUid uuid.UUID) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid})
	if err := r.Delete(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to delete permissions: %v", err)
		return errors.Wrapf(err, "failed to delete permissions for user %s", userUid)
	}

	log.Printf("successfully delete permissions")
	return nil
}
