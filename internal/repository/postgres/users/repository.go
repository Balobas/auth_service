package repositoryUsers

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type UsersRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *UsersRepository {
	return &UsersRepository{
		repositoryPostgres.New(client),
	}
}
