package verification

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *VerificationRepository) CreateVerification(ctx context.Context, verification entity.Verification) error {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(verification)

	if err := r.Create(ctx, verificationRow); err != nil {
		return errors.Wrapf(err, "failed to create verification for user with uid %s", verification.UserUid)
	}
	return nil
}

func (r *VerificationRepository) GetUserVerification(ctx context.Context, userUid uuid.UUID) (entity.Verification, error) {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(entity.Verification{UserUid: userUid})

	if err := r.GetOne(ctx, verificationRow, verificationRow.ConditionUserUidEqual()); err != nil {
		return entity.Verification{}, errors.Wrapf(err, "failed to get verification for user %s", &userUid)
	}

	return verificationRow.ToEntity(), nil
}

func (r *VerificationRepository) UpdateVerification(ctx context.Context, verification entity.Verification) error {
	verififcationRow := pgEntity.NewVerificationRow().FromEntity(verification)

	if err := r.Update(ctx, verififcationRow, verififcationRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to update verification for user %s", verification.UserUid)
	}

	return nil
}

func (r *VerificationRepository) GetVerificationsInStatus(ctx context.Context, status entity.VerificationStatus, limit uint64) ([]entity.Verification, error) {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(entity.Verification{Status: status})
	rows := pgEntity.NewVerificationRows()
	if err := r.GetWithLimit(ctx, verificationRow, rows, verificationRow.ConditionsStatusEqual(), limit, 0); err != nil {
		return nil, errors.Wrapf(err, "failed to get verifications in status %s", status)
	}
	return rows.ToEntities(), nil
}

func (r *VerificationRepository) DeleteVerification(ctx context.Context, userUid uuid.UUID) error {
	verificationRow := pgEntity.NewVerificationRow().FromEntity(entity.Verification{UserUid: userUid})

	if err := r.Delete(ctx, verificationRow, verificationRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete verification for user %s", &userUid)
	}
	return nil
}
