package config

import (
	"time"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/config"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `config:"loglevel"`

	Coeff float64 `config:"coeff"`

	Timeout time.Duration `config:"timeout"`
}

func Init() (Config, error) {
	return config.Load[Config](func() {
		viper.SetDefault("loglevel", "info")
		viper.SetDefault("timeout", "30s")
	})
}
