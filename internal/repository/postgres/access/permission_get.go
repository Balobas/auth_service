package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetPermission(ctx context.Context, key string) (entity.Permission, bool, error) {

	row := pgEntity.NewPermissionRow().FromEntity(entity.Permission{Key: key})

	if err := r.GetOne(ctx, row, row.ConditionUidEqual()); err != nil {
		log.Printf("accessRepository.GetPermission: failed to get permission %s: %v", key, err)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Permission{}, false, nil
		}
		return entity.Permission{}, false, errors.WithStack(err)
	}
	return row.ToEntity(), true, nil
}
