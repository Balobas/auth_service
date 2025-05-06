package workerPermissionsRemover

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	config        Config
	ucPermissions UcPermissions
}

func New(
	config Config,
	ucPermissions UcPermissions,
) *Worker {
	return &Worker{
		config:        config,
		ucPermissions: ucPermissions,
	}
}

func (w *Worker) Run(ctx context.Context) {
	log.Printf("start permissions remover worker\n")
	timer := time.NewTimer(w.config.RemovePermissionsInterval())
	defer timer.Stop() // потестить

	for {
		select {
		case <-ctx.Done():
			log.Printf("permissions remover worker stopped. ctx done")
			return
		case <-timer.C:
			select {
			case <-ctx.Done():
				log.Printf("permissions remover worker stopped. ctx done")
				return
			default:
			}

			w.removePermissionFromUsers(ctx, w.config.UsersLimitOnRemovePermission())
			timer.Reset(w.config.RemovePermissionsInterval())
		}
	}
}

func (w *Worker) removePermissionFromUsers(ctx context.Context, usersLimit int64) {
	perm, isPermDeletingListEmpty, err := w.ucPermissions.GetRandomPermissionFromDeletingList(ctx)
	if err != nil {
		log.Printf("permissionsRemoverWorker: failed to get random permission from deleting list: %v", err)
		return
	}
	if isPermDeletingListEmpty {
		return
	}

	log.Printf("permissionsRemoverWorker: try to remove permission %s from users", perm)

	anyUsersAffected, err := w.ucPermissions.RemovePermissionFromUsers(ctx, perm, usersLimit)
	if err != nil {
		log.Printf("permissionsRemoverWorker: failed to remove permission %s from users: %v", perm, err)
		return
	}

	if anyUsersAffected {
		return
	}

	log.Printf("permissionsRemoverWorker: permission %s was removed from all users having it", perm)

	if err := w.ucPermissions.RemovePermissionFromDeletingList(ctx, perm); err != nil {
		log.Printf("permissionsRemoverWorker: failed to remove permission from deleting list: %v", err)
	}
}
