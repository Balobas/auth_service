package useCaseUsers

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/validations"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) GetUserByUid(ctx context.Context, uid uuid.UUID) (entity.User, bool, error) {
	return uc.usersRepo.GetUserByUid(ctx, uid)
}

func (uc *UseCaseUsers) GetUserByEmail(ctx context.Context, email string) (entity.User, bool, error) {
	if err := validations.ValidateEmail(email); err != nil {
		return entity.User{}, false, errors.WithStack(err)
	}
	return uc.usersRepo.GetByEmail(ctx, email)
}
