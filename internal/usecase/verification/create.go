package useCaseVerification

import (
	"context"
	"crypto/rand"
	"fmt"
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

	if verification.Status == entity.VerificationStatusWaiting {
		verification.Status = entity.VerificationStatusCreated
		verification.UpdatedAt = time.Now()

		if err := uc.verificationRepository.UpdateVerification(ctx, verification); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	for {
		ticker := time.NewTicker(30 * time.Millisecond)

		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "timeout reached")
		case <-ticker.C:
			_, isFound, err := uc.verificationRepository.GetVerificationByToken(ctx, verification.Token)
			if err != nil {
				ticker.Stop()
				return errors.WithStack(err)
			}
			if !isFound {
				defer ticker.Stop()
				if err := uc.verificationRepository.CreateVerification(ctx, verification); err != nil {
					return errors.WithStack(err)
				}
				return nil
			}
		}
	}
}

func randomToken(length int64) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}
