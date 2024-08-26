package repositoryPostgres

import (
	"github.com/balobas/auth_service/internal/client"
)

type Repository struct {
	dbc client.ClientDB
}

type RepositoryConfig interface {
	DSN() string
}

func New(client client.ClientDB) *Repository {
	return &Repository{
		dbc: client,
	}
}

func (r *Repository) db() client.DB {
	return r.dbc.DB()
}
