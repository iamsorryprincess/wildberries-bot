package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSONParser[TConfig any] struct{}

func NewJSONParser[TConfig any]() *JSONParser[TConfig] {
	return &JSONParser[TConfig]{}
}

func (p *JSONParser[TConfig]) Parse(reader io.Reader) (TConfig, error) {
	var cfg TConfig

	if err := json.NewDecoder(reader).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("json parser decode: %w", err)
	}

	return cfg, nil
}
