package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) UpdatePermission(ctx context.Context, perm entity.Permission) error {
	row := pgEntity.NewPermissionRow().FromEntity(perm)

	if err := r.Update(ctx, row, row.ConditionUidEqual()); err != nil {
		log.Printf("accessRepository.UpdatePermission: failed to update permission (key %s): %v", perm.Key, err)
		return errors.WithStack(err)
	}

	return nil
}
