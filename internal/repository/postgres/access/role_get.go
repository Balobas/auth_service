package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *AccessRepository) GetRole(ctx context.Context, role string) (entity.Role, bool, error) {
	log.Printf("accessRepository.GetRole: role %s", role)

	if len(role) == 0 {
		return entity.Role{}, false, errors.New("empty role")
	}

	stmt := "SELECT role, description FROM roles WHERE role=$1"

	row := r.QueryRow(ctx, stmt, role)

	res := entity.Role{}

	if err := row.Scan(&res.Role, &res.Description); err != nil {
		log.Printf("accessRepository.GetRole: failed to get role %s: %v", role, err)

		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Role{}, false, nil
		}
		return entity.Role{}, false, errors.WithStack(err)
	}

	return res, true, nil
}
