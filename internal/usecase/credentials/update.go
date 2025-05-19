package useCaseCredentials

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) Update(ctx context.Context, userUid uuid.UUID, password string) error {
	log.Printf("usecaseCredentials.Update: user %s", userUid)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("usecaseCredentials.Update: failed to get hash from password: %v", err)
		return errors.Wrapf(err, "failed to get hash from password")
	}

	creds := entity.UserCredentials{
		UserUid:      userUid,
		PasswordHash: passwordHash,
	}

	if err := uc.credsRepo.UpdateCredentials(ctx, creds); err != nil {
		log.Printf("usecaseCredentials.Update: failed to update user %s creds: %v", userUid, err)
		return errors.WithStack(err)
	}
	return nil
}
