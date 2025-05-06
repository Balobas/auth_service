package repositoryPermissions

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
