package repositoryPostgres

import "github.com/balobas/auth_service/internal/client"

type BasePgRepository struct {
	dbc client.ClientDB
}

func New(client client.ClientDB) *BasePgRepository {
	return &BasePgRepository{
		dbc: client,
	}
}

func (r *BasePgRepository) DB() client.DB {
	return r.dbc.DB()
}
