package repositoryAccess

import (
	"context"
	"fmt"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetResourcesPermissions(ctx context.Context, limit int64, offset int64) ([]entity.ResourcePermissions, error) {
	log.Printf("accessRepository.GetResourcesPermissions: limit %d, offset %d", limit, offset)

	stmt := "SELECT uri, method, permission FROM resources_permissions "

	args := make([]interface{}, 0, 2)
	idx := 1

	if offset != 0 {
		stmt += fmt.Sprintf("OFFSET $%d", idx)
		idx++
		args = append(args, offset)
	}

	if limit != 0 {
		stmt += fmt.Sprintf("LIMIT $%d", idx)
		idx++
		args = append(args, limit)
	}

	rows, err := r.Query(ctx, stmt, args...)
	if err != nil {
		log.Printf("accessRepository.GetResourcesPermissions: failed to get: %v", err)
		return nil, errors.WithStack(err)
	}

	var res []entity.ResourcePermissions

	for rows.Next() {
		rp := entity.ResourcePermissions{}

		if err := rows.Scan(&rp.URI, &rp.Method, &rp.Permission); err != nil {
			log.Printf("accessRepository.GetResourcesPermissions: failed to scan: %v", err)
			return nil, errors.WithStack(err)
		}
		res = append(res, rp)
	}

	return res, nil
}
