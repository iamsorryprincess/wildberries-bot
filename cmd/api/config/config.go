package config

import "github.com/spf13/viper"

type Config struct {
	LogLevel string `config:"loglevel"`
}

func (c Config) SetDefaults() {
	viper.SetDefault("loglevel", "info")
}
