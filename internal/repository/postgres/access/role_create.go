package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (r *AccessRepository) CreateRole(ctx context.Context, role entity.Role) error {
	log.Printf("accessRepository.CreateRole: role %s, desc %s", role.Role, role.Description)

	if len(role.Role) == 0 {
		return errors.New("empty role")
	}

	stmt := `INSERT INTO roles(role, description) VALUES ($1, $2)`

	_, err := r.Exec(ctx, stmt, role.Role, role.Description)
	if err != nil {
		log.Printf("accessRepository.CreateRole: failed to create role %s: %v", role.Role, err)
		return errors.WithStack(err)
	}
	return nil
}
