package repositoryAccess

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetRoles(ctx context.Context, rolePattern string, limit int64, offset int64) ([]entity.Role, error) {
	log.Printf("accessRepository.GetRoles: rolePattern %s", rolePattern)

	stmt := strings.Builder{}
	stmt.WriteString("SELECT role, description FROM roles ")

	var args []interface{}
	idx := 1

	if len(rolePattern) != 0 {
		stmt.WriteString(fmt.Sprintf("WHERE lower(role) LIKE $%d ", idx))
		idx++
		args = append(args, "%"+strings.ToLower(rolePattern)+"%")
	}

	if offset != 0 {
		stmt.WriteString(fmt.Sprintf("OFFSET $%d", idx))
		idx++
		args = append(args, offset)
	}

	if limit != 0 {
		stmt.WriteString(fmt.Sprintf("LIMIT $%d ", idx))
		idx++
		args = append(args, limit)
	}

	rows, err := r.Query(ctx, stmt.String(), args...)
	if err != nil {
		log.Printf("accessRepository.GetRoles: failed to get roles  (pattern %s): %v", rolePattern, err)
		return nil, errors.WithStack(err)
	}

	res := []entity.Role{}

	for rows.Next() {
		role := entity.Role{}
		if err := rows.Scan(&role.Role, &role.Description); err != nil {
			log.Printf("accessRepository.GetRoles: failed to scan roles: %v", err)
			return nil, errors.WithStack(err)
		}

		res = append(res, role)
	}

	return res, nil
}
