package useCaseVerification

import (
	dbManager "github.com/balobas/sport_city_common/managers/database"
)

type UseCaseVerification struct {
	cfg Config

	verificationRepository VerificationRepository
	usersRepository        UsersRepository
	dbm                    *dbManager.Manager
}

func New(
	cfg Config,
	verificationRepo VerificationRepository,
	usersRepo UsersRepository,
	dbm *dbManager.Manager,
) *UseCaseVerification {
	return &UseCaseVerification{
		cfg:                    cfg,
		verificationRepository: verificationRepo,
		usersRepository:        usersRepo,
		dbm:                    dbm,
	}
}
