package useCaseConfig

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func (uc *UseCaseConfig) Update(ctx context.Context, cfg map[string]json.RawMessage) error {
	if err := uc.config.Validate(cfg); err != nil {
		return errors.Wrap(err, "invalid fields")
	}

	if err := uc.configRepo.UpdateConfig(ctx, cfg); err != nil {
		return errors.Wrap(err, "failed to update config")
	}

	if err := uc.config.LoadFromMap(cfg); err != nil {
		return errors.Wrap(err, "failed to update config")
	}

	return nil
}
