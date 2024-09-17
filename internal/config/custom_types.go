package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if len(str) == 0 {
		return nil
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	if duration < 0 {
		return fmt.Errorf("duration is negative")
	}
	*d = Duration{duration}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Duration < 0 {
		return []byte("null"), nil
	}
	str := d.Duration.String()
	data, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}
