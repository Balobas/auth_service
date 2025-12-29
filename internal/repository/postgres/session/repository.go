package sessionRepository

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type SessionRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *SessionRepository {
	return &SessionRepository{
		repositoryPostgres.New(client),
	}
}
