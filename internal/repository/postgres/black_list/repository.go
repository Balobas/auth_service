package blacklistRepository

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type BlacklistRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *BlacklistRepository {
	return &BlacklistRepository{
		repositoryPostgres.New(client),
	}
}
