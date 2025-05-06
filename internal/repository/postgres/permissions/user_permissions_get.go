package repositoryPermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) GetUserPermissions(ctx context.Context, userUid uuid.UUID) ([]entity.UserPermission, error) {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid})

	if err := r.GetOne(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to get permissions: %v", err)
		return nil, errors.Wrapf(err, "failed to get user %s permissions", userUid)
	}
	usr := &entity.User{}
	permissionsRow.ToEntity(usr)

	log.Printf("successfully get permissions")
	return usr.Permissions, nil
}
