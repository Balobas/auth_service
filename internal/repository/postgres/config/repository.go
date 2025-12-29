package config

import (
	DBclient "github.com/balobas/sport_city_common/clients/database"
)

type ConfigRepository struct {
	client DBclient.ClientDB
}

func New(client DBclient.ClientDB) *ConfigRepository {
	return &ConfigRepository{
		client: client,
	}
}
