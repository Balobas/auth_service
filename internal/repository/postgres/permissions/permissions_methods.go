package permissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) CreateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid, Permissions: perms})
	if err := r.Create(ctx, permissionsRow); err != nil {
		log.Printf("failed to create permissions: %v", err)
		return errors.Wrapf(err, "failed to create permissions for user %s", userUid)
	}

	log.Printf("successfully create permissions")
	return nil
}

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

func (r *PermissionsRepository) UpdateUserPermissions(ctx context.Context, userUid uuid.UUID, perms []entity.UserPermission) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid, Permissions: perms})
	if err := r.Update(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to update permissions %v", err)
		return errors.Wrapf(err, "failed to update permissions for user %s", userUid)
	}

	log.Printf("successfully update permissions")
	return nil
}

func (r *PermissionsRepository) DeleteUserPermissions(ctx context.Context, userUid uuid.UUID) error {
	permissionsRow := pgEntity.NewUserPermissionsRow().FromEntity(entity.User{Uid: userUid})
	if err := r.Delete(ctx, permissionsRow, permissionsRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to delete permissions: %v", err)
		return errors.Wrapf(err, "failed to delete permissions for user %s", userUid)
	}

	log.Printf("successfully delete permissions")
	return nil
}
