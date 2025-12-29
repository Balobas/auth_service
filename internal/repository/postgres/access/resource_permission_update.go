package repositoryAccess

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *AccessRepository) UpdateResourcePermission(ctx context.Context, uri string, method string, permission string) error {
	log.Printf("accessRepository.UpdateResourcePermission: uri %s, method %s, permission %s", uri, method, permission)

	stmt := "UPDATE resources_permissions SET permission=$1 WHERE uri=$2 AND method=$3"

	if _, err := r.Exec(ctx, stmt, permission, uri, method); err != nil {
		log.Printf("accessRepository.UpdateResourcePermission: failed to update permission (%s) for resource %s method %s: %v", permission, uri, method, err)
		return errors.WithStack(err)
	}

	return nil
}
