package outboxRepository

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type OutboxRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *OutboxRepository {
	return &OutboxRepository{
		repositoryPostgres.New(client),
	}
}
