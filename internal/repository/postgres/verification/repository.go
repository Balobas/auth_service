package verification

import (
	"github.com/balobas/auth_service/internal/client"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
)

type VerificationRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client client.ClientDB) *VerificationRepository {
	return &VerificationRepository{
		BasePgRepository: repositoryPostgres.New(client),
	}
}
