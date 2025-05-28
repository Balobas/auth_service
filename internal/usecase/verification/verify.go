package useCaseVerification

import (
	"context"
	"log"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
)

func (uc *UseCaseVerification) Verify(ctx context.Context, token string) error {
	log.Printf("usecaseVerification.Verify: token %s", token)

	if len(token) == 0 {
		log.Printf("usecaseVerification.Verify: empty token")
		return errors.Wrap(serviceErrors.ErrBadRequest, "empty token")
	}

	if err := uc.txManager.NewPgTransaction().Execute(ctx, func(ctx context.Context) error {
		verification, isFound, err := uc.verificationRepository.GetVerificationByToken(ctx, token)
		if err != nil {
			log.Printf("usecaseVerification.Verify: failed to get verification by token %s: %v", token, err)
			return errors.WithStack(err)
		}

		if !isFound {
			log.Printf("usecaseVerification.Verify: verification not found by token %s", token)
			return errors.Wrap(serviceErrors.ErrNotFound, "verification not found by token")
		}

		user, isFound, err := uc.usersRepository.GetUserByUid(ctx, verification.UserUid)
		if err != nil {
			log.Printf("usecaseVerification.Verify: failed to get user by uid %s: %v", verification.UserUid, err)
			return errors.WithStack(err)
		}

		if !isFound {
			log.Printf("usecaseVerification.Verify: user not found by uid %s", verification.UserUid)
			return errors.Wrap(serviceErrors.ErrNotFound, "user")
		}
		user.IsVerified = true

		if err := uc.usersRepository.UpdateUser(ctx, user); err != nil {
			log.Printf("usecaseVerification.Verify: failed to update user %s: %v", verification.UserUid, err)
			return errors.WithStack(err)
		}

		if err := uc.verificationRepository.DeleteVerification(ctx, verification.UserUid); err != nil {
			log.Printf("usecaseVerification.Verify: failed to delete user %s verification: %v", verification.UserUid, err)
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
