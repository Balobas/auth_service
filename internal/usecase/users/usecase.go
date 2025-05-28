package useCaseUsers

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseUsers struct {
	cfg        Config
	usersRepo  UsersRepository
	accessRepo AccessRepository

	ucVerification   UcVerification
	ucCredentials    UcCredentials
	ucOutboxMessages UcOutboxMessages
	jwtManager       JwtManager
	txManager        *transaction.Manager
}

func New(
	cfg Config,
	usersRepo UsersRepository,
	accessRepo AccessRepository,
	ucVerification UcVerification,
	txManager *transaction.Manager,
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
		txManager:        txManager,
		jwtManager:       jwtManager,
	}
}
