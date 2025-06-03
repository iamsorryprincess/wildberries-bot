package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

func LoadYAML[TConfig Config](filename string) (TConfig, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	var cfg TConfig

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	cfg.SetDefaults()

	if err := viper.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "config"
		config.IgnoreUntaggedFields = true
	}); err != nil {
		return cfg, err
	}

	return cfg, nil
}
