package useCaseAuth

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseAuth struct {
	cfg Config

	sessionsRepo SessionsRepository
	permsRepo    PermissionsRepository

	ucUsers       UcUsers
	ucCredentials UcCredentials

	jwtManager JwtManager
	txManager  *transaction.Manager
}

func New(
	cfg Config,
	sessionsRepo SessionsRepository,
	permsRepo PermissionsRepository,
	ucUsers UcUsers,
	ucCreds UcCredentials,
	jwtManager JwtManager,
	txManager *transaction.Manager,
) *UseCaseAuth {
	return &UseCaseAuth{
		cfg:           cfg,
		sessionsRepo:  sessionsRepo,
		permsRepo:     permsRepo,
		ucUsers:       ucUsers,
		ucCredentials: ucCreds,
		jwtManager:    jwtManager,
		txManager:     txManager,
	}
}
