package devicesRepository

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type DevicesRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *DevicesRepository {
	return &DevicesRepository{BasePgRepository: repositoryPostgres.New(client)}
}
