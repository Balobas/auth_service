package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) DeleteUser(ctx context.Context, userUid uuid.UUID) error {
	log.Printf("usecaseAuth.DeleteUser: user %s", userUid)

	// все связное должно каскадом удалиться

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		if err := uc.usersRepo.DeleteUser(ctx, userUid); err != nil {
			log.Printf("usecaseAuth.DeleteUser: failed to delete user %s: %v", userUid, err)
			return err
		}

		if err := uc.ucOutboxMessages.CreateUserDeletedMessage(ctx, entity.User{Uid: userUid}); err != nil {
			log.Printf("usecaseAuth.DeleteUser: failed to create user %s deleted message: %v", userUid, err)
			return err
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
