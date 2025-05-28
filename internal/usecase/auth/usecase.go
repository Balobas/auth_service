package useCaseAuth

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseAuth struct {
	cfg Config

	sessionsRepo SessionsRepository
	accessRepo   AccessRepository

	ucUsers       UcUsers
	ucCredentials UcCredentials

	jwtManager JwtManager
	txManager  *transaction.Manager
}

func New(
	cfg Config,
	sessionsRepo SessionsRepository,
	accessRepo AccessRepository,
	ucUsers UcUsers,
	ucCreds UcCredentials,
	jwtManager JwtManager,
	txManager *transaction.Manager,
) *UseCaseAuth {
	return &UseCaseAuth{
		cfg:           cfg,
		sessionsRepo:  sessionsRepo,
		accessRepo:    accessRepo,
		ucUsers:       ucUsers,
		ucCredentials: ucCreds,
		jwtManager:    jwtManager,
		txManager:     txManager,
	}
}
