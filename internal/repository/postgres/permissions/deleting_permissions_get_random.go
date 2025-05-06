package repositoryPermissions

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (r *PermissionsRepository) GetRandomPermissionFromDeletingList(ctx context.Context) (perm string, isListEmpty bool, err error) {
	log.Printf("permissionsRepository.GetRandomPermissionFromDeletingList")

	stmt := "select key from deleting_permissions tablesample system(0.2) limit 1"
	row := r.DB().QueryRow(ctx, stmt)
	if err := row.Scan(&perm); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", true, nil
		}

		log.Printf("permissionsRepository.GetRandomPermissionFromDeletingList: failed to get: %v", err)
		return "", false, errors.WithStack(err)
	}
	return perm, false, nil
}
