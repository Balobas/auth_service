package verification

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
)

type VerificationRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *VerificationRepository {
	return &VerificationRepository{
		BasePgRepository: repositoryPostgres.New(client),
	}
}
