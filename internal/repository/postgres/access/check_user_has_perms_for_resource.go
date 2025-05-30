package repositoryAccess

import (
	"context"
	"log"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *AccessRepository) IsUserHasPermissionsForResource(ctx context.Context, userUid uuid.UUID, uri string, method string) (bool, error) {
	log.Printf("repositoryAccess.IsUserHasPermissionsForResource: userUid: %s, uri: %s, method: %s", userUid, uri, method)

	stmt := `with perms as (
		select distinct on (rp.permission) rp.permission from user_roles ur inner join roles_permissions rp on ur.role = rp.role where ur.user_uid = $1
		)
		select true as res from resources_permissions where uri = $2 and method = $3 and (permission in (select permission from perms) or permission = 'none');
	`

	args := []interface{}{
		pgtype.UUID{
			Bytes:  userUid,
			Status: pgtype.Present,
		},
		uri,
		method,
	}

	row := r.DB().QueryRow(ctx, stmt, args...)

	var res bool

	err := row.Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		log.Printf("repositoryAccess.IsUserHasPermissionsForResource: failed to scan row: %v", err)
		return false, err
	}

	return true, nil
}
