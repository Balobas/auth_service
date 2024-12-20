package useCaseUsers

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
)

func (uc *UseCaseUsers) GetAdmins(ctx context.Context, token string) ([]entity.User, error) {
	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("GetAdmins: failed to parse token\n")
		return []entity.User{}, errors.WithStack(err)
	}

	_, isFound, err := uc.getAdminUser(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("GetAdmins: failed to parse token\n")
		return []entity.User{}, errors.WithStack(err)
	}

	if !isFound {
		log.Printf("GetAdmins: failed to find user\n")
		return []entity.User{}, errors.WithStack(err)
	}

	users, err := uc.usersRepo.GetAdminUsers(ctx)
	if err != nil {
		log.Printf("GetAdmins: failed to get admin users\n")
		return []entity.User{}, errors.WithStack(err)
	}

	return users, nil
}
