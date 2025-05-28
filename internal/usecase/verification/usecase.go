package useCaseVerification

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseVerification struct {
	cfg Config

	verificationRepository VerificationRepository
	usersRepository        UsersRepository
	txManager              *transaction.Manager
}

func New(
	cfg Config,
	verificationRepo VerificationRepository,
	usersRepo UsersRepository,
	txManager *transaction.Manager,
) *UseCaseVerification {
	return &UseCaseVerification{
		cfg:                    cfg,
		verificationRepository: verificationRepo,
		usersRepository:        usersRepo,
		txManager:              txManager,
	}
}
