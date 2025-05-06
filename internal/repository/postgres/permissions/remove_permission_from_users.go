package repositoryPermissions

import (
	"context"
	"log"

	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *PermissionsRepository) RemovePermissionFromUsers(ctx context.Context, key string, limitUsers int64) ([]uuid.UUID, error) {

	stmt := `update user_permissions set permissions = array_remove(permissions, $1) 
		where user_uid in (select user_uid from user_permissions where $2=any(permissions) limit $3)
		returning user_uid`

	args := make([]interface{}, 3)
	args[0] = key
	args[1] = key
	args[2] = limitUsers

	rows, err := r.DB().Query(ctx, stmt, args...)
	if err != nil {
		log.Printf("repositoryPermissions.RemovePermissionFromUsers: failed to remove permission %s from users: %v", key, err)
		return nil, errors.WithStack(err)
	}

	uidRows := pgEntity.NewUUIDRows()
	if err := uidRows.ScanAll(rows); err != nil {
		log.Printf("repositoryPermissions.RemovePermissionFromUsers: failed to scan users uids: %v", err)
		return nil, errors.WithStack(err)
	}
	return uidRows.ToEntity(), nil
}
