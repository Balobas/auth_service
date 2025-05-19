package useCaseVerification

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

//TODO: засунуть в транзацию, чтобы не оставалось висящих верификаций

func (uc *UseCaseVerification) Verify(ctx context.Context, token string) error {
	log.Printf("usecaseVerification.Verify: token %s", token)

	verification, isFound, err := uc.verificationRepository.GetVerificationByToken(ctx, token)
	if err != nil {
		log.Printf("usecaseVerification.Verify: failed to get verification by token %s: %v", token, err)
		return errors.WithStack(err)
	}

	if !isFound {
		return errors.Wrap(serviceErrors.ErrNotFound, "verification not found by token")
	}

	// TODO: отрефакторить изменение прав
	// если пользователь у которого были различные права меняет email, то все его права сбросятся,
	// а после верификации проставится только base
	if err := uc.permissionsRepository.UpdateUserPermissions(
		ctx,
		verification.UserUid,
		[]entity.UserPermission{entity.UserPermissionBase},
	); err != nil {
		log.Printf("usecaseVerification.Verify: failed to update user %s permissions: %v", verification.UserUid, err)
		return errors.WithStack(err)
	}

	if err := uc.verificationRepository.DeleteVerification(ctx, verification.UserUid); err != nil {
		log.Printf("usecaseVerification.Verify: failed to delete user %s verification: %v", verification.UserUid, err)
		return errors.WithStack(err)
	}

	return nil
}
