package repositoryPostgres

import (
	"context"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/entity/contract"
	"github.com/pkg/errors"
)

type BasePgRepository struct {
	dbc client.ClientDB
}

func New(client client.ClientDB) *BasePgRepository {
	return &BasePgRepository{
		dbc: client,
	}
}

func (r *BasePgRepository) DB() client.DB {
	return r.dbc.DB()
}

func (r *BasePgRepository) WithTx(ctx context.Context, f func(ctx context.Context) error) (err error) {
	if !r.DB().HasTxInCtx(ctx) {
		var (
			tx         contract.Transaction
			beginTxErr error
		)
		ctx, tx, beginTxErr = r.DB().BeginTxWithContext(ctx)
		if beginTxErr != nil {
			return errors.Wrap(beginTxErr, "failed to start transaction")
		}

		defer func() {
			err = HandleTxEnd(ctx, tx, err)
		}()
	}

	return f(ctx)
}

func (r *BasePgRepository) BeginTxWithContext(ctx context.Context) (context.Context, contract.Transaction, error) {
	return r.dbc.DB().BeginTxWithContext(ctx)
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
