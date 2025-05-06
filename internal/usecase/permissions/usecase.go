package useCasePermissions

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCasePermissions struct {
	permsRepository    PermissionsRepository
	usersRepository    UsersRepository
	sessionsRepository SessionsRepository
	txManager          *transaction.Manager
}
