package useCaseAuth

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// TODO: не удалять сессию из БД.
func (uc *UseCaseAuth) Logout(ctx context.Context, userUid uuid.UUID) error {
	log.Printf("usecaseAuth.Logout: user %s", userUid)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecaseAuth.Logout: empty user uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	if err := uc.sessionsRepo.DeleteSessionByUserUid(ctx, userUid); err != nil {
		log.Printf("usecaseAuth.Logout: failed to delete session for user %s: %v", userUid, err)
		return errors.WithStack(err)
	}
	return nil
}
