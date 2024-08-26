package repositoryUsers

import (
	"context"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/entity/contract"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
	"github.com/pkg/errors"
)

type UsersRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *UsersRepository {
	return &UsersRepository{
		repositoryPostgres.New(client),
	}
}

func HandleTxEnd(ctx context.Context, tx contract.Transaction, err error) error {
	if err == nil {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return errors.Wrap(commitErr, "failed to commit transaction")
		}
		return nil
	}

	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		return errors.Wrap(rollbackErr, "failed to rollback transaction")
	}
	return nil
}
