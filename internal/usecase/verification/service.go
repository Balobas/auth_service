package useCaseVerification

type UseCaseVerification struct {
	cfg Config

	verificationRepository VerificationRepository
	permissionsRepository  PermissionsRepository
}

func New(
	cfg Config,
	verificationRepo VerificationRepository,
	permissionsRepo PermissionsRepository,
) *UseCaseVerification {
	return &UseCaseVerification{
		cfg:                    cfg,
		verificationRepository: verificationRepo,
		permissionsRepository:  permissionsRepo,
	}
}
