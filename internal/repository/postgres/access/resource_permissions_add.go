package repositoryAccess

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *AccessRepository) AddResourcePermission(ctx context.Context, uri string, method string, permission string) error {
	log.Printf("accessRepository.AddResourcePermission: uri %s, method %s, permission %s", uri, method, permission)

	stmt := "INSERT INTO resources_permissions (uri, method, permission) VALUES ($1, $2, $3)"

	if _, err := r.DB().Exec(ctx, stmt, uri, method, permission); err != nil {
		log.Printf("accessRepository.AddResourcePermission: failed to add permission %s to resource (method %s, uri %s): %v", permission, method, uri, err)
		return errors.WithStack(err)
	}

	return nil
}
