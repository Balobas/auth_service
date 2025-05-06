package useCasePermissions

import (
	"context"
	"log"
)

func (uc *UseCasePermissions) RemovePermissionFromDeletingList(ctx context.Context, permKey string) error {
	log.Printf("usecasePermissions.RemovePermissionFromUsers: perm %s", permKey)

	return uc.permsRepository.RemovePermissionFromDeletingList(ctx, permKey)
}
