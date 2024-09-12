package useCaseAuth

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseAuth) Logout(ctx context.Context, userUid uuid.UUID) error {
	if err := uc.sessionsRepo.DeleteSessionByUserUid(ctx, userUid); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
