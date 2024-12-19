package useCaseUsers

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseUsers struct {
	usersRepo UsersRepository
	permsRepo PermissionsRepository

	ucVerification   UcVerification
	ucCredentials    UcCredentials
	ucOutboxMessages UcOutboxMessages
	jwtManager       JwtManager
	txManager        *transaction.Manager
}

func New(
	usersRepo UsersRepository,
	permsRepo PermissionsRepository,
	ucVerification UcVerification,
	txManager *transaction.Manager,
	ucCreds UcCredentials,
	jwtManager JwtManager,
	ucOutboxMessages UcOutboxMessages,
) *UseCaseUsers {
	return &UseCaseUsers{
		usersRepo:        usersRepo,
		permsRepo:        permsRepo,
		ucVerification:   ucVerification,
		ucCredentials:    ucCreds,
		ucOutboxMessages: ucOutboxMessages,
		txManager:        txManager,
		jwtManager:       jwtManager,
	}
}
