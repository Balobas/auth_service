package repositoryAccess

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (r *AccessRepository) DeleteRole(ctx context.Context, role string) error {
	log.Printf("accessRepository.DeleteRole: role %s", role)

	if len(role) == 0 {
		return errors.New("empty role")
	}

	stmt := `DELETE FROM roles WHERE role=$1`

	_, err := r.Exec(ctx, stmt, role)
	if err != nil {
		log.Printf("accessRepository.DeleteRole: failed to delete role %s: %v", role, err)
		return errors.WithStack(err)
	}
	return nil
}
