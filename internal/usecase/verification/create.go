package useCaseVerification

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	common "github.com/balobas/sport_city_common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseVerification) CreateVerification(ctx context.Context, userUid uuid.UUID, email string) error {
	log.Printf("usecaseVerification.CreateVerification: user %s email %s", userUid, email)

	if uuid.Equal(userUid, uuid.UUID{}) {
		log.Printf("usecaseVerification.CreateVerification: empty user uid")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user uid")
	}
	if len(email) == 0 {
		log.Printf("usecaseVerification.CreateVerification: empty user email")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty user email")
	}

	verification := entity.Verification{
		UserUid:   userUid,
		Token:     randomToken(uc.cfg.VerificationTokenLen()),
		Email:     email,
		Status:    entity.VerificationStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.dbm.ExecuteTx(ctx, common.Serializable, func(ctx context.Context) error {

		oldVerification, isFound, err := uc.verificationRepository.GetUserVerification(ctx, userUid)
		if err != nil {
			log.Printf("usecaseVerification.CreateVerification: failed to get user %s old verification: %v", userUid, err)
			return err
		}
		if oldVerification.Email == email {
			log.Printf("usecaseVerification.CreateVerification: verification for user with email %s already exists", email)
			return errors.Wrap(serviceErrors.ErrAlreadyExists, "verification for user email already exists")
		}

		if isFound {
			if err := uc.verificationRepository.DeleteVerification(ctx, userUid); err != nil {
				log.Printf("usecaseVerification.CreateVerification: failed to delete old verification for user %s: %v", userUid, err)
				return err
			}
		}
		if err := uc.verificationRepository.CreateVerification(ctx, verification); err != nil {
			log.Printf("usecaseVerification.CreateVerification: failed to create verification for user %s: %v", userUid, err)
			return err
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func randomToken(length int64) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
