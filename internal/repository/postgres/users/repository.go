package repositoryUsers

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type UsersRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *UsersRepository {
	return &UsersRepository{
		repositoryPostgres.New(client),
	}
}
