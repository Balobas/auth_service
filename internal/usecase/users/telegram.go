package useCaseUsers

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// GetUserByTelegramID возвращает пользователя по Telegram id (found=false, если нет).
func (uc *UseCaseUsers) GetUserByTelegramID(ctx context.Context, telegramID int64) (entity.User, bool, error) {
	return uc.usersRepo.GetByTelegramID(ctx, telegramID)
}

// RegisterTelegramUser заводит пользователя, пришедшего из Telegram: без пароля и
// email-верификации, сразу verified. name/nickname (first_name/username) уходят в
// событие UserRegistered для создания профиля в users_service.
func (uc *UseCaseUsers) RegisterTelegramUser(ctx context.Context, telegramID int64, name, nickname string) (entity.User, error) {
	user := entity.User{
		Uid:        uuid.NewV4(),
		TelegramID: telegramID,
		IsVerified: true,
		CreatedAt:  time.Now(),
	}

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		if err := uc.usersRepo.CreateUser(ctx, user); err != nil {
			return err
		}
		if err := uc.accessRepo.AddRoleToUser(ctx, user.Uid, entity.UserRoleUser); err != nil {
			return err
		}
		return uc.ucOutboxMessages.CreateUserRegisteredMessage(ctx, user, name, nickname)
	}); err != nil {
		return entity.User{}, err
	}

	user.Roles = []string{entity.UserRoleUser}
	return user, nil
}

// SetTelegramID привязывает Telegram к существующему аккаунту (линковка). email и
// прочие поля сохраняются (грузим юзера целиком перед апдейтом).
func (uc *UseCaseUsers) SetTelegramID(ctx context.Context, userUid uuid.UUID, telegramID int64) error {
	return uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {
		user, found, err := uc.usersRepo.GetUserByUid(ctx, userUid)
		if err != nil {
			return err
		}
		if !found {
			return errors.Wrap(serviceErrors.ErrNotFound, "user for telegram linking not found")
		}
		user.TelegramID = telegramID
		user.UpdatedAt = time.Now()
		return uc.usersRepo.UpdateUser(ctx, user)
	})
}
