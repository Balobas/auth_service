package useCaseAuth

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// TODO: не удалять сессию из БД.
func (uc *UseCaseAuth) Logout(ctx context.Context, userUid uuid.UUID) error {
	if err := uc.sessionsRepo.DeleteSessionByUserUid(ctx, userUid); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
