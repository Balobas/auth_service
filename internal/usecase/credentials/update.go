package useCaseCredentials

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) Update(ctx context.Context, userUid uuid.UUID, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrapf(err, "failed to get hash from password")
	}

	creds := entity.UserCredentials{
		UserUid:      userUid,
		PasswordHash: passwordHash,
	}

	if err := uc.credsRepo.UpdateCredentials(ctx, creds); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
