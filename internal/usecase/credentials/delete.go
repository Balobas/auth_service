package useCaseCredentials

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseCredentials) Delete(ctx context.Context, userUid uuid.UUID) error {
	if err := uc.credsRepo.DeleteByUserUid(ctx, userUid); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
