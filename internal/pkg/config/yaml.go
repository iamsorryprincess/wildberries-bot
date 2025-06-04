package config

import (
	"bytes"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type YamlParser[TConfig any] struct{}

func NewYamlParser[TConfig any]() *YamlParser[TConfig] {
	return &YamlParser[TConfig]{}
}

func (p *YamlParser[TConfig]) Parse(data []byte) (TConfig, error) {
	viper.SetConfigType("yaml")

	var cfg TConfig
	if err := viper.ReadConfig(bytes.NewReader(data)); err != nil {
		return cfg, fmt.Errorf("config yaml parser read: %w", err)
	}

	if err := viper.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "config"
		config.IgnoreUntaggedFields = true
	}); err != nil {
		return cfg, fmt.Errorf("config yaml parser unmarshal: %w", err)
	}

	return cfg, nil
}
