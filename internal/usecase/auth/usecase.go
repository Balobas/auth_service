package useCaseAuth

import (
	dbManager "github.com/balobas/sport_city_common/managers/database"
)

type UseCaseAuth struct {
	cfg Config

	sessionsRepo SessionsRepository
	accessRepo   AccessRepository
	servicesRepo ServicesRepository

	ucUsers       UcUsers
	ucCredentials UcCredentials
	ucDevices     UcDevices

	jwtManager JwtManager
	dbm        *dbManager.Manager
}

func New(
	cfg Config,
	sessionsRepo SessionsRepository,
	accessRepo AccessRepository,
	servicesRepo ServicesRepository,
	ucUsers UcUsers,
	ucCreds UcCredentials,
	ucDevices UcDevices,
	jwtManager JwtManager,
	dbm *dbManager.Manager,
) *UseCaseAuth {
	return &UseCaseAuth{
		cfg:           cfg,
		sessionsRepo:  sessionsRepo,
		accessRepo:    accessRepo,
		servicesRepo:  servicesRepo,
		ucUsers:       ucUsers,
		ucCredentials: ucCreds,
		ucDevices:     ucDevices,
		jwtManager:    jwtManager,
		dbm:           dbm,
	}
}
