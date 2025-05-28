package repositoryAccess

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type AccessRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *AccessRepository {
	return &AccessRepository{
		repositoryPostgres.New(client),
	}
}
