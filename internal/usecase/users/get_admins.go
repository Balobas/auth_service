package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

// TODO: рефактор: вынести валидацию токена в delivery, логики получения это не должно касаться!

func (uc *UseCaseUsers) GetAdmins(ctx context.Context) ([]entity.User, error) {
	log.Printf("usecaseAuth.GetAdmins: ")

	users, err := uc.usersRepo.GetAdminUsers(ctx)
	if err != nil {
		log.Printf("usecaseAuth.GetAdmins: failed to get admin users: %v", err)
		return []entity.User{}, errors.WithStack(err)
	}
	if len(users) == 0 {
		log.Printf("usecaseAuth.GetAdmins: no admins found")
		return nil, errors.Wrap(serviceErrors.ErrNotFound, "admins")
	}

	return users, nil
}
