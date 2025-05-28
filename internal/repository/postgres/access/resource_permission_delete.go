package repositoryAccess

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *AccessRepository) DeleteResourcePermission(ctx context.Context, uri string, method string, permission string) error {
	log.Printf("accessRepository.DeleteResourcePermission: uri %s, permission %s", uri, permission)

	stmt := "DELETE FROM resources_permissions WHERE WHERE uri=$1 AND method=$2"

	if _, err := r.DB().Exec(ctx, stmt, uri, method); err != nil {
		log.Printf("accessRepository.DeleteResourcePermission: failed to delete resource (method %s, uri %s) permissions: %v", method, uri, err)
		return errors.WithStack(err)
	}

	return nil
}
