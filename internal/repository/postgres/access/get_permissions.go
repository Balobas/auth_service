package repositoryAccess

import (
	"context"
	"log"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetPermissions(ctx context.Context, keyPattern string) ([]entity.Permission, error) {
	log.Printf("accessRepository.GetPermissions: key pattern %s", keyPattern)

	row := pgEntity.NewPermissionRow()
	rows := pgEntity.NewPermissionsRows()

	var cond squirrel.Sqlizer
	if len(keyPattern) != 0 {
		cond = squirrel.Like{
			"lower(key)": "%" + strings.ToLower(keyPattern) + "%",
		}
	}

	if err := r.GetSome(ctx, row, rows, cond); err != nil {
		log.Printf("accessRepository.GetPermissions: failed to get permissions (keyPattern %s): %v", keyPattern, err)
		return nil, errors.WithStack(err)
	}
	return rows.ToEntity(), nil
}
