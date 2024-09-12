package useCaseCredentials

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) Validate(ctx context.Context, userUid uuid.UUID, password string) error {
	creds, isFound, err := uc.credsRepo.GetByUserUid(ctx, userUid)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isFound {
		return errors.Errorf("credentials for user %s not found", userUid)
	}

	return bcrypt.CompareHashAndPassword(creds.PasswordHash, []byte(password))
}
