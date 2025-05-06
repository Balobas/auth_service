package useCasePermissions

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCasePermissions) AddPermissionToUser(ctx context.Context, userUid uuid.UUID, permKey string) error {
	log.Printf("usecasePermissions.AddPermissionToUser: user %s permission %s", userUid, permKey)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecasePermissions.AddPermissionToUser: empty user uid")
		return errors.New("empty user uid")
	}
	if len(permKey) == 0 {
		log.Printf("usecasePermissions.AddPermissionToUser: empty permission")
		return errors.New("empty permission")
	}

	err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {

		_, isFound, err := uc.usersRepository.GetUserByUid(ctx, userUid)
		if err != nil {
			log.Printf("usecasePermissions.AddPermissionToUser: failed to get user %s: %v", userUid, err)
			return err
		}
		if !isFound {
			log.Printf("usecasePermissions.AddPermissionToUser: user %s not found", userUid)
			return errors.New("user not found")
		}

		_, isFound, err = uc.permsRepository.GetPermission(ctx, permKey)
		if err != nil {
			log.Printf("usecasePermissions.AddPermissionToUser: failed to get permission %s: %v", permKey, err)
			return err
		}
		if !isFound {
			log.Printf("usecasePermissions.AddPermissionToUser: permission %s not found", permKey)
			return errors.New("permission not found")
		}

		perms, err := uc.permsRepository.GetUserPermissions(ctx, userUid)
		if err != nil {
			log.Printf("usecasePermissions.AddPermissionToUser: failed to get user %s permissions: %v", userUid, err)
			return err
		}

		permsMap := make(map[entity.UserPermission]struct{})
		for i := 0; i < len(perms); i++ {
			permsMap[perms[i]] = struct{}{}
		}

		if _, ok := permsMap[entity.UserPermission(permKey)]; ok {
			log.Printf("usecasePermissions.AddPermissionToUser: user %s already have permission %s", userUid, permKey)
			return errors.Errorf("user already have permission %s", permKey)
		}

		perms = append(perms, entity.UserPermission(permKey))
		if err := uc.permsRepository.UpdateUserPermissions(ctx, userUid, perms); err != nil {
			log.Printf("usecasePermissions.AddPermissionToUser: failed to update user %s permissions: %v", userUid, err)
			return err
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
