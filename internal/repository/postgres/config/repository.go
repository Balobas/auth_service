package config

import (
	"github.com/balobas/auth_service/internal/client"
)

type ConfigRepository struct {
	client client.ClientDB
}

func New(client client.ClientDB) *ConfigRepository {
	return &ConfigRepository{
		client: client,
	}
}
