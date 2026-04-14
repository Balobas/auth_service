package devicesRepository

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type Repository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *Repository {
	return &Repository{
		repositoryPostgres.New(client),
	}
}
