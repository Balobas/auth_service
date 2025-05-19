package useCaseVerification

import (
	"context"
	"log"

	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseVerification) DeleteUserVerification(ctx context.Context, userUid uuid.UUID) error {
	log.Printf("usecaseVerification.DeleteUserVerification: user %s", userUid)

	return uc.verificationRepository.DeleteVerification(ctx, userUid)
}
