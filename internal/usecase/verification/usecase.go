package useCaseVerification

import "github.com/balobas/auth_service/internal/manager/transaction"

type UseCaseVerification struct {
	cfg Config

	verificationRepository VerificationRepository
	permissionsRepository  PermissionsRepository
	txManager              *transaction.Manager
}

func New(
	cfg Config,
	verificationRepo VerificationRepository,
	permissionsRepo PermissionsRepository,
	txManager *transaction.Manager,
) *UseCaseVerification {
	return &UseCaseVerification{
		cfg:                    cfg,
		verificationRepository: verificationRepo,
		permissionsRepository:  permissionsRepo,
		txManager:              txManager,
	}
}
