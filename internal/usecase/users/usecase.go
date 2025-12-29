package useCaseUsers

import (
	dbManager "github.com/balobas/sport_city_common/managers/database"
)

type UseCaseUsers struct {
	cfg        Config
	usersRepo  UsersRepository
	accessRepo AccessRepository

	ucVerification   UcVerification
	ucCredentials    UcCredentials
	ucOutboxMessages UcOutboxMessages
	jwtManager       JwtManager
	dbm              *dbManager.Manager
}

func New(
	cfg Config,
	usersRepo UsersRepository,
	accessRepo AccessRepository,
	ucVerification UcVerification,
	dbm *dbManager.Manager,
	ucCreds UcCredentials,
	jwtManager JwtManager,
	ucOutboxMessages UcOutboxMessages,
) *UseCaseUsers {
	return &UseCaseUsers{
		cfg:              cfg,
		usersRepo:        usersRepo,
		accessRepo:       accessRepo,
		ucVerification:   ucVerification,
		ucCredentials:    ucCreds,
		ucOutboxMessages: ucOutboxMessages,
		dbm:              dbm,
		jwtManager:       jwtManager,
	}
}
