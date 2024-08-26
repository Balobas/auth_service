package transaction

import (
	"context"
	"testing"

	"github.com/balobas/auth_service/internal/client/pg"
)

func TestTx(t *testing.T) {
	pgClient, err := pg.NewClient(context.Background(), "")
	if err != nil {
		t.Fail()
	}

	txManager := NewTxManager()

	tx := txManager.NewTransaction(pgClient.DB())

	if err := tx.Execute(context.Background(), func(ctx context.Context) error {
		_, err := pgClient.DB().Exec(ctx, "Select hui from huis")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		t.Fail()
	}
}
