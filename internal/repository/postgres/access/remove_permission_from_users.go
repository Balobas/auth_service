package repositoryAccess

import (
	"context"
	"log"

	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Deprecated
func (r *AccessRepository) RemovePermissionFromUsers(ctx context.Context, key string, limitUsers int64) ([]uuid.UUID, error) {
	if len(key) == 0 {
		log.Printf("accessRepository.RemovePermissionFromUsers: empty permission")
		return nil, errors.New("empty permission")
	}

	stmt := `update user_permissions set permissions = array_remove(permissions, $1) 
		where user_uid in (select user_uid from user_permissions where $2=any(permissions)`

	if limitUsers != 0 {
		stmt += "limit $3"
	}
	stmt += ") returning user_uid"

	args := make([]interface{}, 2, 3)
	args[0] = key
	args[1] = key
	if limitUsers != 0 {
		args = append(args, limitUsers)
	}

	rows, err := r.DB().Query(ctx, stmt, args...)
	if err != nil {
		log.Printf("accessRepository.RemovePermissionFromUsers: failed to remove permission %s from users: %v", key, err)
		return nil, errors.WithStack(err)
	}

	uidRows := pgEntity.NewUUIDRows()
	if err := uidRows.ScanAll(rows); err != nil {
		log.Printf("accessRepository.RemovePermissionFromUsers: failed to scan users uids: %v", err)
		return nil, errors.WithStack(err)
	}
	return uidRows.ToEntity(), nil
}
