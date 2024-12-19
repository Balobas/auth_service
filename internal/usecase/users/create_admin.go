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
		log.Printf("verifyAuth: failed to parse token\n")
		return uuid.UUID{}, errors.WithStack(err)
	}

	if tokenInfo.Role != string(entity.UserRoleAdmin) {
		log.Printf("role from token is invalid: cannot create admin with this role")
		return uuid.UUID{}, fmt.Errorf("role is invalid")
	}

	uid, err := uc.Register(ctx, user, password)
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return uid, nil
}
