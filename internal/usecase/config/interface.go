package useCaseConfig

import (
	"context"
	"encoding/json"
)

type Config interface {
	LoadFromMap(config map[string]json.RawMessage) error
	ToMap() map[string]json.RawMessage
	Validate(cfg map[string]json.RawMessage) error
}

type ConfigRepository interface {
	UpdateConfig(ctx context.Context, cfg map[string]json.RawMessage) error
	GetConfig(ctx context.Context) (map[string]json.RawMessage, error)
}
