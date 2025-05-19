package useCaseCredentials

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) Validate(ctx context.Context, userUid uuid.UUID, password string) error {
	log.Printf("usecaseCredentials.Validate: user %s", userUid)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecaseCredentials.Validate: empty user uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}

	if len(password) == 0 {
		log.Printf("usecaseCredentials.Validate: empty password")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty password")
	}

	creds, isFound, err := uc.credsRepo.GetByUserUid(ctx, userUid)
	if err != nil {
		log.Printf("usecaseCredentials.Validate: failed to get user %s creds: %v", userUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseCredentials.Validate: creds for user %s not found", userUid)
		// internal error, креды должны быть для юзера который есть в бд
		return errors.Errorf("credentials for user %s not found", userUid)
	}

	return bcrypt.CompareHashAndPassword(creds.PasswordHash, []byte(password))
}
