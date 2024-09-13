package useCaseConfig

import (
	"context"
	"encoding/json"
)

func (uc *UseCaseConfig) Get(ctx context.Context) map[string]json.RawMessage {
	return uc.config.ToMap()
}
