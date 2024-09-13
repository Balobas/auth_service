package useCaseConfig

import (
	"context"

	"github.com/pkg/errors"
)

func (uc *UseCaseConfig) InitFromDB(ctx context.Context) error {
	cfgMap, err := uc.configRepo.GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	if err := uc.config.LoadFromMap(cfgMap); err != nil {
		return errors.Wrap(err, "failed to init config")
	}
	return nil
}
