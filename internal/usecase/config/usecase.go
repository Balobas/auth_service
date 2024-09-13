package useCaseConfig

type UseCaseConfig struct {
	config     Config
	configRepo ConfigRepository
}

func New(cfg Config, cfgRepo ConfigRepository) *UseCaseConfig {
	return &UseCaseConfig{
		config:     cfg,
		configRepo: cfgRepo,
	}
}
