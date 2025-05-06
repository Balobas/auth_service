package useCasePermissions

import (
	"context"
	"log"
)

func (uc *UseCasePermissions) GetRandomPermissionFromDeletingList(ctx context.Context) (perm string, isListEmpty bool, err error) {
	log.Printf("usecasePermissions.GetRandomPermissionFromDeletingList")
	return uc.permsRepository.GetRandomPermissionFromDeletingList(ctx)
}
