package repositoryAccess

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetResourcePermission(ctx context.Context, uri string, method string) (string, bool, error) {
	log.Printf("accessRepository.GetResourcePermission: method %s uri %s", method, uri)

	stmt := "SELECT permission FROM resources_permissions WHERE uri=$1 AND method = $2"

	row := r.DB().QueryRow(ctx, stmt, uri, method)

	var perm string
	if err := row.Scan(&perm); err != nil {
		log.Printf("accessRepository.GetResourcePermission: failed to get method %s uri %s permissions: %v", method, uri, err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, errors.WithStack(err)
	}
	return perm, true, nil
}
