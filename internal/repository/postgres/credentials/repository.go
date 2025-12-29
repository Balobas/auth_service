package credentials

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type CredentialsRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *CredentialsRepository {
	return &CredentialsRepository{
		repositoryPostgres.New(client),
	}
}
