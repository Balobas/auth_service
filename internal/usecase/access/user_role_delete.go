package useCaseAccess

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAccess) RemoveRoleFromUser(ctx context.Context, userUid uuid.UUID, role string) error {
	log.Printf("usecasePemrissions.RemoveRoleFromUser: user %s role %s", userUid, role)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecasePemrissions.RemoveRoleFromUser: empty user uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	if len(role) == 0 {
		log.Printf("usecasePemrissions.RemoveRoleFromUser: empty role")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty role")
	}

	if err := uc.accessRepository.DeleteRoleFromUser(ctx, userUid, role); err != nil {
		log.Printf("usecasePemrissions.RemoveRoleFromUser: failed to remove role %s from user %s: %v", role, userUid, err)
		return err
	}

	return nil
}
