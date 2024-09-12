package permissions

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) CreateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid, Permissions: perms})
	if err := r.Create(ctx, permissionsRow); err != nil {
		return errors.Wrapf(err, "failed to create permissions for user %s", userUid)
	}
	return nil
}

func (r *PermissionsRepository) GetUserPermissions(ctx context.Context, userUid uuid.UUID) ([]entity.UserPermission, error) {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid})

	if err := r.GetOne(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		return nil, errors.Wrapf(err, "failed to get user %s permissions", userUid)
	}
	return pgEntity.NewUserRow().ToEntity().Permissions, nil
}

func (r *PermissionsRepository) UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid, Permissions: perms})
	if err := r.Update(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to update permissions for user %s", userUid)
	}
	return nil
}

func (r *PermissionsRepository) DeleteUserPermissions(ctx context.Context, userUid uuid.UUID) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid})
	if err := r.Delete(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete permissions for user %s", userUid)
	}
	return nil
}
