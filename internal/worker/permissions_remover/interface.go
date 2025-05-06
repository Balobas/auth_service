package workerPermissionsRemover

import (
	"context"
	"time"
)

type Config interface {
	UsersLimitOnRemovePermission() int64
	RemovePermissionsInterval() time.Duration
}

type UcPermissions interface {
	RemovePermissionFromUsers(ctx context.Context, permKey string, limitUsers int64) (anyUsersAffected bool, err error)
	GetRandomPermissionFromDeletingList(ctx context.Context) (perm string, isListEmpty bool, err error)
	RemovePermissionFromDeletingList(ctx context.Context, permKey string) error
}
