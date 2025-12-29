package repositoryAccess

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type AccessRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *AccessRepository {
	return &AccessRepository{
		repositoryPostgres.New(client),
	}
}
