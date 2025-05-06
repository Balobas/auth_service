package useCasePermissions

import (
	"context"
	"log"

	"github.com/pkg/errors"
)

func (uc *UseCasePermissions) RemovePermissionFromUsers(ctx context.Context, permKey string, limitUsers int64) (bool, error) {
	log.Printf("usecasePermissions.RemovePermissionFromUsers: permission %s", permKey)

	if len(permKey) == 0 {
		log.Printf("usecasePermissions.RemovePermissionFromUsers: empty permission")
		return false, errors.New("empty permission")
	}

	var anyUsersAffected bool

	err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {

		affectedUsers, err := uc.permsRepository.RemovePermissionFromUsers(ctx, permKey, limitUsers)
		if err != nil {
			log.Printf("usecasePermissions.RemovePermissionFromUsers: failed to remove permission %s from users: %v", permKey, err)
			return err
		}
		if len(affectedUsers) == 0 {
			log.Printf("usecasePermissions.RemovePermissionFromUsers: users with permission %s not found", permKey)
			return nil
		}

		anyUsersAffected = true
		// Удаляем сессии, чтобы пользователи вызвали рефреш токена
		if err := uc.sessionsRepository.DeleteSessionsByUsersUids(ctx, affectedUsers); err != nil {
			log.Printf("usecasePermissions.RemovePermissionFromUsers: failed to delete users sessions: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		return false, errors.WithStack(err)
	}

	return anyUsersAffected, nil
}
