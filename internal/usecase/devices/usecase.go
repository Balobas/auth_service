package ucDevices

import dbManager "github.com/balobas/sport_city_common/managers/database"

type UseCase struct {
	dbm         *dbManager.Manager
	devicesRepo DevicesRepository
}

func New(
	dbm *dbManager.Manager,
	devicesRepo DevicesRepository,
) *UseCase {
	return &UseCase{
		dbm:         dbm,
		devicesRepo: devicesRepo,
	}
}
