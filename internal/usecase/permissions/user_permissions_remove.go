package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCasePermissions) RemoveUserPermission(ctx context.Context, userUid uuid.UUID, permKey string) error {
	log.Printf("usecasePermissions.RemoveUserPermission: user %s permission %s", userUid, permKey)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecasePermissions.RemoveUserPermission: empty user uid")
		return errors.New("empty user uid")
	}
	if len(permKey) == 0 {
		log.Printf("usecasePermissions.RemoveUserPermission: empty permission")
		return errors.New("empty permission")
	}

	err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		_, isFound, err := uc.usersRepository.GetUserByUid(ctx, userUid)
		if err != nil {
			log.Printf("usecasePermissions.RemoveUserPermission: failed to get user %s: %v", userUid, err)
			return err
		}
		if !isFound {
			log.Printf("usecasePermissions.RemoveUserPermission: user %s not found", userUid)
			return errors.New("user not found")
		}

		perms, err := uc.permsRepository.GetUserPermissions(ctx, userUid)
		if err != nil {
			log.Printf("usecasePermissions.RemoveUserPermission: failed to get user %s permissions: %v", userUid, err)
			return err
		}

		perm := entity.UserPermission(permKey)
		userHasRemovingPerm := false
		removingIdx := -1

		for idx, p := range perms {
			if p == perm {
				userHasRemovingPerm = true
				removingIdx = idx
			}
		}

		if !userHasRemovingPerm {
			log.Printf("usecasePermissions.RemoveUserPermission: user %s hasnt permission %s", userUid, permKey)
			return errors.Errorf("user hasnt permission %s", permKey)
		}

		perms[removingIdx] = perms[len(perms)-1]
		perms = perms[:len(perms)-1]

		if err := uc.permsRepository.UpdateUserPermissions(ctx, userUid, perms); err != nil {
			log.Printf("usecasePermissions.RemoveUserPermission: failed to remove user %s permission %s: %v", userUid, permKey, err)
			return err
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
