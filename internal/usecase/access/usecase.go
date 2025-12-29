package useCaseAccess

import (
	dbManager "github.com/balobas/sport_city_common/managers/database"
)

type UseCaseAccess struct {
	accessRepository AccessRepository
	usersRepository  UsersRepository
	dbm              *dbManager.Manager
}

func New(
	accessRepository AccessRepository,
	usersRepository UsersRepository,
	dbm *dbManager.Manager,
) *UseCaseAccess {
	return &UseCaseAccess{
		accessRepository: accessRepository,
		usersRepository:  usersRepository,
		dbm:              dbm,
	}
}
