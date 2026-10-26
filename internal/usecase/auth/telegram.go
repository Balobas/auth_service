package useCaseAuth

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/balobas/auth_service/pkg/telegram"
	common "github.com/balobas/sport_city_common"
	uuid "github.com/satori/go.uuid"
)

// Максимальный возраст initData (anti-replay). Мини-апп отдаёт свежую при каждом открытии.
const telegramInitDataMaxAge = 24 * time.Hour

// AuthWithTelegram — вход через Telegram Mini App по initData. Если telegram ещё не
// привязан: при CreateIfMissing=true заводит аккаунт, иначе возвращает
// ErrTelegramUserNotFound (клиент показывает выбор «есть аккаунт?»).
func (uc *UseCaseAuth) AuthWithTelegram(ctx context.Context, params entity.TelegramAuthParams) (access string, refresh string, created bool, err error) {
	tgUser, err := telegram.ValidateInitData(uc.cfg.TelegramBotToken(), params.InitData, telegramInitDataMaxAge)
	if err != nil {
		return "", "", false, err
	}

	user, found, err := uc.ucUsers.GetUserByTelegramID(ctx, tgUser.ID)
	if err != nil {
		return "", "", false, err
	}
	if !found {
		if !params.CreateIfMissing {
			return "", "", false, serviceErrors.ErrTelegramUserNotFound
		}
		user, err = uc.ucUsers.RegisterTelegramUser(ctx, tgUser.ID, tgUser.FirstName, tgUser.Username)
		if err != nil {
			return "", "", false, err
		}
		created = true
	}

	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		u, e := uc.ensureRoles(ctx, user)
		if e != nil {
			return e
		}
		access, refresh, e = uc.issueUserSession(ctx, u, params.Device, time.Now())
		return e
	}); err != nil {
		return "", "", false, err
	}

	return access, refresh, created, nil
}

// LinkTelegram — привязка Telegram к существующему аккаунту (вход по email+паролю).
// Связь строго 1:1.
func (uc *UseCaseAuth) LinkTelegram(ctx context.Context, params entity.TelegramLinkParams) (access string, refresh string, err error) {
	tgUser, err := telegram.ValidateInitData(uc.cfg.TelegramBotToken(), params.InitData, telegramInitDataMaxAge)
	if err != nil {
		return "", "", err
	}

	user, err := uc.ucUsers.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return "", "", err
	}
	if err := uc.ucCredentials.Validate(ctx, user.Uid, params.Password); err != nil {
		return "", "", err
	}

	// Правила 1:1.
	existingByTg, foundByTg, err := uc.ucUsers.GetUserByTelegramID(ctx, tgUser.ID)
	if err != nil {
		return "", "", err
	}
	if foundByTg && !uuid.Equal(existingByTg.Uid, user.Uid) {
		return "", "", serviceErrors.ErrTelegramAlreadyLinked
	}
	if user.TelegramID != 0 && user.TelegramID != tgUser.ID {
		return "", "", serviceErrors.ErrAccountAlreadyHasTelegram
	}

	if user.TelegramID != tgUser.ID {
		if err := uc.ucUsers.SetTelegramID(ctx, user.Uid, tgUser.ID); err != nil {
			return "", "", err
		}
		user.TelegramID = tgUser.ID
	}

	if err := uc.dbm.ExecuteTx(ctx, common.ReadCommitted, func(ctx context.Context) error {
		u, e := uc.ensureRoles(ctx, user)
		if e != nil {
			return e
		}
		access, refresh, e = uc.issueUserSession(ctx, u, params.Device, time.Now())
		return e
	}); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// ensureRoles догружает роли пользователя, если они ещё не заполнены.
func (uc *UseCaseAuth) ensureRoles(ctx context.Context, user entity.User) (entity.User, error) {
	if len(user.Roles) != 0 {
		return user, nil
	}
	roles, err := uc.accessRepo.GetUserRoles(ctx, user.Uid)
	if err != nil {
		return entity.User{}, err
	}
	user.Roles = entity.RolesToStrings(roles)
	return user, nil
}

// issueUserSession — общая логика выдачи токенов: переиспользует сессию для
// (user, device) или создаёт новую. Выполняется в транзакции вызывающего.
func (uc *UseCaseAuth) issueUserSession(ctx context.Context, user entity.User, deviceData entity.LoginDeviceData, loginTime time.Time) (string, string, error) {
	device, err := entity.NewUserDevice(user.Uid, deviceData)
	if err != nil {
		return "", "", err
	}

	savedSession, isFound, err := uc.sessionsRepo.GetSessionByMaintainerAndDevice(ctx, user.Uid, device.Uid)
	if err != nil {
		return "", "", err
	}
	if isFound {
		return uc.handleUserLoginWithExistingSession(user, savedSession)
	}

	if err := uc.ucDevices.HandleLoginFromDevice(ctx, device, loginTime); err != nil {
		return "", "", err
	}

	session := entity.NewUserSession(uuid.NewV4(), user.Uid, device.Uid, loginTime)
	if err := uc.sessionsRepo.CreateSession(ctx, session); err != nil {
		return "", "", err
	}

	tokenInfo := entity.NewUserTokenInfo(user, device.Uid, session.Uid, loginTime)
	return uc.jwtManager.NewUserTokens(tokenInfo, uc.cfg.AccessJwtTTL(), uc.cfg.RefreshJwtTTL())
}
