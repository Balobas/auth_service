package repositoryUsers

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *UsersRepository) CreateVerification(ctx context.Context, verification entity.Verification) error {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(verification)

	if err := r.Create(ctx, verificationRow); err != nil {
		return errors.Wrapf(err, "failed to create verification for user with uid %s", verification.UserUid)
	}
	return nil
}

func (r *UsersRepository) GetUserVerification(ctx context.Context, userUid uuid.UUID) (entity.Verification, error) {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(entity.Verification{UserUid: userUid})

	if err := r.GetOne(ctx, verificationRow, verificationRow.ConditionUserUidEqual()); err != nil {
		return entity.Verification{}, errors.Wrapf(err, "failed to get verification for user %s", &userUid)
	}

	return verificationRow.ToEntity(), nil
}

func (r *UsersRepository) DeleteVerification(ctx context.Context, userUid uuid.UUID) error {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(entity.Verification{UserUid: userUid})

	if err := r.Delete(ctx, verificationRow, verificationRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete verification for user %s", &userUid)
	}
	return nil
}
