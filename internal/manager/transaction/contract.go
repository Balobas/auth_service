package transaction

import (
	"context"

	"github.com/balobas/auth_service/internal/entity/contract"
)

type Transactor interface {
	BeginTxWithContext(ctx context.Context) (context.Context, contract.Transaction, error)
	HasTxInCtx(ctx context.Context) bool
}
