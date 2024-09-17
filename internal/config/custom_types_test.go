package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCustomDuration(t *testing.T) {
	data := `"1m"`

	d := new(Duration)

	fmt.Println(json.Unmarshal([]byte(data), &d))
	fmt.Println(*d)
}
