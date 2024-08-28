package sessionRepository

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type SessionRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *SessionRepository {
	return &SessionRepository{
		repositoryPostgres.New(client),
	}
}
