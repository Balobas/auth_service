package useCaseUsers

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) DeleteUser(ctx context.Context, userUid uuid.UUID) error {
	// все связное должно какадом удалиться
	if err := uc.usersRepo.DeleteUser(ctx, userUid); err !=nil {
		return errors.WithStack(err)
	}
	return nil
}