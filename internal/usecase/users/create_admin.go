package useCaseUsers

import (
	"context"
	"fmt"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) CreateAdmin(ctx context.Context, user entity.User, password string, token string) (uuid.UUID, error) {
	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("CreateAdmin: failed to parse token\n")
		return uuid.UUID{}, errors.WithStack(err)
	}

	creatingUser, isFound, err := uc.usersRepo.GetByEmail(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("CreateAdmin: failed to get user\n")
		return uuid.UUID{}, fmt.Errorf("failed to get user")
	}

	if !isFound {
		log.Printf("CreateAdmin: failed to find user\n")
		return uuid.UUID{}, fmt.Errorf("failed to find user")
	}

	if creatingUser.Role != entity.UserRoleAdmin || creatingUser.Role != entity.UserRole(tokenInfo.Role) {
		log.Printf("role from token is invalid: cannot create admin with this role")
		return uuid.UUID{}, fmt.Errorf("role is invalid")
	}

	user.Role = entity.UserRoleAdmin

	log.Printf("CreateAdmin: registering user\n")

	uid, err := uc.Register(ctx, user, password)
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return uid, nil
}
