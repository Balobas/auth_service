package credentials

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type CredentialsRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *CredentialsRepository {
	return &CredentialsRepository{
		repositoryPostgres.New(client),
	}
}
