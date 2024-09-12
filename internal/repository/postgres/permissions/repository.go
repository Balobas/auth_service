package permissions

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type PermissionsRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *PermissionsRepository {
	return &PermissionsRepository{
		repositoryPostgres.New(client),
	}
}
