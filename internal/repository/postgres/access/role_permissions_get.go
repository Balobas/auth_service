package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetRolePermissions(ctx context.Context, role string) ([]entity.Permission, error) {
	log.Printf("accessRepository.GetRolePermissions: role %s", role)

	stmt := "SELECT permission, p.description FROM roles_permissions rp INNER JOIN permissions p on rp.permission=p.key WHERE role=$1"

	rows, err := r.DB().Query(ctx, stmt, role)
	if err != nil {
		log.Printf("accessRepository.GetRolePermissions: failed to get role %s permissions: %v", role, err)
		return nil, errors.WithStack(err)
	}

	var perms []entity.Permission

	for rows.Next() {
		var perm entity.Permission
		if err := rows.Scan(&perm.Key, &perm.Description); err != nil {
			log.Printf("accessRepository.GetRolePermissions: failed to scan role %s permissions: %v", role, err)
			return nil, errors.WithStack(err)
		}
		perms = append(perms, perm)
	}
	return perms, nil
}
