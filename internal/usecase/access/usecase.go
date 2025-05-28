package useCaseAccess

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseAccess struct {
	accessRepository AccessRepository
	usersRepository  UsersRepository
	txManager        *transaction.Manager
}

func New(
	accessRepository AccessRepository,
	usersRepository UsersRepository,
	txManager *transaction.Manager,
) *UseCaseAccess {
	return &UseCaseAccess{
		accessRepository: accessRepository,
		usersRepository:  usersRepository,
		txManager:        txManager,
	}
}
