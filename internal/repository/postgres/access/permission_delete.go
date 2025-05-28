package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) DeletePermission(ctx context.Context, key string) error {

	row := pgEntity.NewPermissionRow().FromEntity(entity.Permission{Key: key})

	if err := r.Delete(ctx, row, row.ConditionUidEqual()); err != nil {
		log.Printf("accessRepository.DeletePermission: failed to delete permission %s: %v", key, err)
		return errors.WithStack(err)
	}
	return nil
}
