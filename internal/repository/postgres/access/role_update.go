package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) UpdateRole(ctx context.Context, role entity.Role) error {
	log.Printf("accessRepository.UpdateRole: role %s, desc %s", role.Role, role.Description)

	if len(role.Role) == 0 {
		return errors.New("empty role")
	}

	stmt := `UPDATE roles SET description=$1 WHERE role=$2`

	_, err := r.Exec(ctx, stmt, role.Description, role.Role)
	if err != nil {
		log.Printf("accessRepository.UpdateRole: failed to update role %s: %v", role.Role, err)
		return errors.WithStack(err)
	}
	return nil
}
