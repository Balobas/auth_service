package useCaseAuth

type UseCaseAuth struct {
	cfg Config

	sessionsRepo SessionsRepository
	permsRepo    PermissionsRepository

	ucUsers       UcUsers
	ucCredentials UcCredentials

	jwtManager JwtManager
}

func New(
	cfg Config,
	sessionsRepo SessionsRepository,
	permsRepo PermissionsRepository,
	ucUsers UcUsers,
	ucCreds UcCredentials,
	jwtManager JwtManager,
) *UseCaseAuth {
	return &UseCaseAuth{
		cfg:           cfg,
		sessionsRepo:  sessionsRepo,
		permsRepo:     permsRepo,
		ucUsers:       ucUsers,
		ucCredentials: ucCreds,
		jwtManager:    jwtManager,
	}
}
