package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) DeleteUser(ctx context.Context, userUid uuid.UUID) error {
	// все связное должно каскадом удалиться

	if err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		if err := uc.usersRepo.DeleteUser(ctx, userUid); err != nil {
			return errors.WithStack(err)
		}

		if err := uc.ucOutboxMessages.CreateUserDeletedMessage(ctx, entity.User{Uid: userUid}); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		log.Printf("failed to delete user %s: %v", userUid, err)
		return errors.WithStack(err)
	}

	return nil
}
