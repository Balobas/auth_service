package useCaseVerification

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseVerification) CreateVerification(ctx context.Context, userUid uuid.UUID, email string) error {
	verification := entity.Verification{
		UserUid:   userUid,
		Token:     randomToken(uc.cfg.VerificationTokenLen()),
		Email:     email,
		Status:    entity.VerificationStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx := uc.txManager.NewPgTransaction()
	if err := tx.Execute(ctx, func(ctx context.Context) error {
		
		oldVerification, isFound, err := uc.verificationRepository.GetUserVerification(ctx, userUid)
		if err != nil {
			return err
		}
		if oldVerification.Email == email {
			log.Printf("verification for user with email %s already exists", email)
			return nil
		}

		if isFound {
			if err := uc.verificationRepository.DeleteVerification(ctx, userUid); err != nil {
				return err
			}
		}
		if err := uc.verificationRepository.CreateVerification(ctx, verification); err != nil {
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
