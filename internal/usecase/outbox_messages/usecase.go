package useCaseOutboxMessages

type UseCaseOutboxMessages struct {
	cfg              Config
	riverClient  RiverClient
}

func New(
	cfg Config,
	riverClient RiverClient,
) *UseCaseOutboxMessages {
	return &UseCaseOutboxMessages{
		cfg:              cfg,
		riverClient: riverClient,
	}
}
