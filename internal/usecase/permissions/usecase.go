package useCasePermissions

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCasePermissions struct {
	permsRepository PermissionsRepository
	usersRepository UsersRepository
	txManager       *transaction.Manager
}

func New(
	permsRepository PermissionsRepository,
	usersRepository UsersRepository,
	txManager *transaction.Manager,
) *UseCasePermissions {
	return &UseCasePermissions{
		permsRepository: permsRepository,
		usersRepository: usersRepository,
		txManager:       txManager,
	}
}
