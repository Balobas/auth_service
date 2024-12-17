package useCaseOutboxMessages

type UseCaseOutboxMessages struct {
	cfg              Config
	outboxRepository OutboxRepository
}

func New(
	cfg Config,
	outboxRepository OutboxRepository,
) *UseCaseOutboxMessages {
	return &UseCaseOutboxMessages{
		cfg:              cfg,
		outboxRepository: outboxRepository,
	}
}
