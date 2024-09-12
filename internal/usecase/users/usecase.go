package useCaseUsers

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseUsers struct {
	usersRepo UsersRepository
	permsRepo PermissionsRepository

	ucVerification UcVerification
	ucCredentials  UcCredentials
	txManager      transaction.Manager
}

func New(
	usersRepo UsersRepository,
	permsRepo PermissionsRepository,
	ucVerification UcVerification,
	txManager transaction.Manager,
	ucCreds UcCredentials,
) *UseCaseUsers {
	return &UseCaseUsers{
		usersRepo:      usersRepo,
		permsRepo:      permsRepo,
		ucVerification: ucVerification,
		ucCredentials:  ucCreds,
		txManager:      txManager,
	}
}
