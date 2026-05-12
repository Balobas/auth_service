package useCaseCredentials

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *UseCaseCredentials) CreateServiceCreds(ctx context.Context, serviceUid uuid.UUID, password string) error {
	log.Printf("usecaseCredentials.CreateServiceCreds: service %s", serviceUid)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("usecaseCredentials.CreateServiceCreds: failed to get hash from password: %v", err)
		return errors.Wrapf(err, "failed to get hash from password")
	}

	creds := entity.ServiceCredentials{
		ServiceUid:   serviceUid,
		PasswordHash: passwordHash,
	}

	if err := uc.credsRepo.CreateServiceCredentials(ctx, creds); err != nil {
		log.Printf("usecaseCredentials.CreateServiceCreds: failed to create (service %s): %v", serviceUid, err)
		return err
	}
	return nil
}

func (uc *UseCaseCredentials) ValidateServiceCreds(ctx context.Context, serviceUid uuid.UUID, password string) error {
	log.Printf("usecaseCredentials.ValidateServiceCreds: service %s", serviceUid)

	if uuid.Equal(serviceUid, uuid.UUID{}) {
		log.Printf("usecaseCredentials.ValidateServiceCreds: empty service uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty service uid")
	}

	if len(password) == 0 {
		log.Printf("usecaseCredentials.ValidateServiceCreds: empty password")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty password")
	}

	creds, isFound, err := uc.credsRepo.GetByServiceUid(ctx, serviceUid)
	if err != nil {
		log.Printf("usecaseCredentials.ValidateServiceCreds: failed to get service %s creds: %v", serviceUid, err)
		return errors.WithStack(err)
	}
	if !isFound {
		log.Printf("usecaseCredentials.ValidateServiceCreds: creds for service %s not found", serviceUid)
		// internal error, креды должны быть для юзера который есть в бд
		return errors.Errorf("credentials for service %s not found", serviceUid)
	}

	return bcrypt.CompareHashAndPassword(creds.PasswordHash, []byte(password))
}
