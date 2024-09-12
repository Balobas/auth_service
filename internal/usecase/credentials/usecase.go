package useCaseCredentials

type UseCaseCredentials struct {
	cfg       Config
	credsRepo CredentialsRepository
}

func New(cfg Config, repo CredentialsRepository) *UseCaseCredentials {
	return &UseCaseCredentials{
		cfg:       cfg,
		credsRepo: repo,
	}
}
