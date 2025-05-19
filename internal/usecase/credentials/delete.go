package useCaseCredentials

import (
	"context"
	"log"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseCredentials) Delete(ctx context.Context, userUid uuid.UUID) error {
	log.Printf("usecaseCredentials.Delete: user %s", userUid)
	if err := uc.credsRepo.DeleteByUserUid(ctx, userUid); err != nil {
		log.Printf("usecaseCredentials.Create: failed to delete user %s creds: %v", userUid, err)
		return errors.WithStack(err)
	}
	return nil
}
