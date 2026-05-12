package useCaseServices

import (
	dbManager "github.com/balobas/sport_city_common/managers/database"
)

type UcServices struct {
	cfg          Config
	servicesRepo ServicesRepository
	accessRepo   AccessRepository
	ucCreds      UcCredentials
	dbm          *dbManager.Manager
}

func New(
	cfg Config,
	servicesRepo ServicesRepository,
	accessRepo AccessRepository,
	ucCreds UcCredentials,
	dbm *dbManager.Manager,
) *UcServices {
	return &UcServices{
		cfg:          cfg,
		servicesRepo: servicesRepo,
		accessRepo:   accessRepo,
		ucCreds:      ucCreds,
		dbm:          dbm,
	}
}
