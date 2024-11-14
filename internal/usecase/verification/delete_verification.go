package useCaseVerification

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseVerification) DeleteUserVerification(ctx context.Context, userUid uuid.UUID) error {
	return uc.verificationRepository.DeleteVerification(ctx, userUid)
}
