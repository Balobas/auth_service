package ucDevices

import dbManager "github.com/balobas/sport_city_common/managers/database"

type UseCase struct {
	dbm         *dbManager.Manager
	devicesRepo DevicesRepository
}

func New(
	devicesRepo DevicesRepository,
) *UseCase {
	return &UseCase{
		devicesRepo: devicesRepo,
	}
}
