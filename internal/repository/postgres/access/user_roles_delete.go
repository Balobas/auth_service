package repositoryAccess

import (
	"context"
	"log"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *AccessRepository) DeleteRoleFromUser(ctx context.Context, userUid uuid.UUID, role string) error {
	log.Printf("accessRepository.DeleteRoleFromUser: user %s role %s", userUid, role)

	stmt := "DELETE FROM user_roles WHERE user_uid=$1 AND role=$2"

	_, err := r.DB().Exec(ctx, stmt, []interface{}{
		pgtype.UUID{
			Bytes:  userUid,
			Status: pgtype.Present,
		},
		role,
	}...)
	if err != nil {
		log.Printf("accessRepository.DeleteRoleFromUser: failed to delete role %s to user %s: %v", role, userUid, err)
		return errors.WithStack(err)
	}

	return nil
}
