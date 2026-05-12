package credentials

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (r *CredentialsRepository) CreateCredentials(ctx context.Context, creds entity.UserCredentials) error {
	credsRow := pgEntity.NewUserCredentialsRow().FromEntity(creds)
	if err := r.Create(ctx, credsRow); err != nil {
		return errors.Wrapf(err, "failed to create credentials for user %s", creds.UserUid)
	}
	return nil
}

func (r *CredentialsRepository) UpdateCredentials(ctx context.Context, creds entity.UserCredentials) error {
	credsRow := pgEntity.NewUserCredentialsRow().FromEntity(creds)
	if err := r.Update(ctx, credsRow, credsRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to create credentials for user %s", creds.UserUid)
	}
	return nil
}

func (r *CredentialsRepository) GetByUserUid(ctx context.Context, userUid uuid.UUID) (entity.UserCredentials, bool, error) {
	credsRow := pgEntity.NewUserCredentialsRow().FromEntity(entity.UserCredentials{UserUid: userUid})

	if err := r.GetOne(ctx, credsRow, credsRow.ConditionUserUidEqual()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserCredentials{}, false, nil
		}
		return entity.UserCredentials{}, false, errors.Wrapf(err, "failed to get credentials for user %s", userUid)
	}

	return credsRow.ToEntity(), true, nil
}

func (r *CredentialsRepository) DeleteByUserUid(ctx context.Context, userUid uuid.UUID) error {
	credsRow := pgEntity.NewUserCredentialsRow().FromEntity(entity.UserCredentials{UserUid: userUid})
	if err := r.Delete(ctx, credsRow, credsRow.ConditionUserUidEqual()); err != nil {
		return errors.Wrapf(err, "failed to delete credentials of user %s", userUid)
	}
	return nil
}

func (r *CredentialsRepository) CreateServiceCredentials(ctx context.Context, creds entity.ServiceCredentials) error {
	credsRow := pgEntity.NewServiceCredentialsRow().FromEntity(creds)
	if err := r.Create(ctx, credsRow); err != nil {
		return errors.Wrapf(err, "failed to create credentials for service %s", creds.ServiceUid)
	}
	return nil
}

func (r *CredentialsRepository) GetByServiceUid(ctx context.Context, serviceUid uuid.UUID) (entity.ServiceCredentials, bool, error) {
	credsRow := pgEntity.NewServiceCredentialsRow()

	if err := r.GetOne(ctx, credsRow, credsRow.ConditionServiceUidEqual(serviceUid)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ServiceCredentials{}, false, nil
		}
		return entity.ServiceCredentials{}, false, errors.Wrapf(err, "failed to get credentials for service %s", serviceUid)
	}

	return credsRow.ToEntity(), true, nil
}
