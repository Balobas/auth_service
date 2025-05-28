package repositoryAccess

import (
	"context"
	"log"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *AccessRepository) AddRoleToUser(ctx context.Context, userUid uuid.UUID, role string) error {
	log.Printf("accessRepository.AddRoleToUser: user %s role %s", userUid, role)

	stmt := "INSERT INTO user_roles (user_uid, role) VALUES ($1, $2)"

	_, err := r.DB().Exec(ctx, stmt, []interface{}{
		pgtype.UUID{
			Bytes:  userUid,
			Status: pgtype.Present,
		},
		role,
	}...)
	if err != nil {
		log.Printf("accessRepository.AddRoleToUser: failed to add role %s to user %s: %v", role, userUid, err)
		return errors.WithStack(err)
	}

	return nil
}
