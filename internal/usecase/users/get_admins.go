package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

// TODO: рефактор: вынести валидацию токена в delivery, логики получения это не должно касаться!

func (uc *UseCaseUsers) GetAdmins(ctx context.Context, token string) ([]entity.User, error) {
	log.Printf("usecaseAuth.GetAdmins: ")

	// вынести в delivery от сюдова
	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("usecaseAuth.GetAdmins: failed to parse token\n")
		return []entity.User{}, errors.Wrap(serviceErrors.ErrBadRequest, "failed to parse token")
	}

	_, isFound, err := uc.getAdminUser(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("usecaseAuth.GetAdmins:  failed to get caller admin by email %s: %v", tokenInfo.Email, err)
		return []entity.User{}, errors.WithStack(err)
	}

	if !isFound {
		log.Printf("usecaseAuth.GetAdmins: caller admin with email %s not found", tokenInfo.Email)
		return []entity.User{}, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "caller is not admin")
	}

	// до сюдова

	users, err := uc.usersRepo.GetAdminUsers(ctx)
	if err != nil {
		log.Printf("usecaseAuth.GetAdmins: failed to get admin users: %v", err)
		return []entity.User{}, errors.WithStack(err)
	}
	if len(users) == 0 {
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "admins")
	}

	return users, nil
}
