package useCaseVerification

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCaseVerification) Verify(ctx context.Context, token string) error {
	verification, isFound, err := uc.verificationRepository.GetVerificationByToken(ctx, token)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isFound {
		return errors.Errorf("verification token %s not found", token)
	}

	if err := uc.permissionsRepository.UpdateUserPermissions(
		ctx,
		verification.UserUid,
		[]entity.UserPermission{entity.UserPermissionBase},
	); err != nil {
		return errors.WithStack(err)
	}

	if err := uc.verificationRepository.DeleteVerification(ctx, verification.UserUid); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
