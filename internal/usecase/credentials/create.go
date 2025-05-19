package useCaseCredentials

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) Create(ctx context.Context, userUid uuid.UUID, password string) error {
	log.Printf("usecaseCredentials.Create: user %s", userUid)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("usecaseCredentials.Create: failed to get hash from password: %v", err)
		return errors.Wrapf(err, "failed to get hash from password")
	}

	creds := entity.UserCredentials{
		UserUid:      userUid,
		PasswordHash: passwordHash,
	}

	if err := uc.credsRepo.CreateCredentials(ctx, creds); err != nil {
		log.Printf("usecaseCredentials.Create: failed to create (user %s): %v", userUid, err)
		return errors.WithStack(err)
	}
	return nil
}
