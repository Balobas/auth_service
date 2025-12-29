package repositoryAccess

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *AccessRepository) GetUserRoles(ctx context.Context, userUid uuid.UUID) ([]entity.Role, error) {
	log.Printf("accessRepository.GetUserRoles: user %s", userUid)

	stmt := "SELECT ur.role, r.description from user_roles ur INNER JOIN roles r on ur.role=r.role  WHERE user_uid=$1"

	rows, err := r.Query(ctx, stmt, pgtype.UUID{
		Bytes:  userUid,
		Valid: true,
	})
	if err != nil {
		log.Printf("accessRepository.GetUserRoles: failed to get user %s roles: %v", userUid, err)
		return nil, errors.WithStack(err)
	}

	var roles []entity.Role

	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.Role, &role.Description); err != nil {
			log.Printf("accessRepository.GetUserRoles: failed to scan user %s roles: %v", userUid, err)
			return nil, errors.WithStack(err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}
