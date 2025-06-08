package config

import (
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/config"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `config:"loglevel"`

	HTTPConfig http.Config `config:"http"`
}

func Init() (Config, error) {
	return config.Load[Config](func() {
		viper.SetDefault("loglevel", "info")

		viper.SetDefault("http.port", "8080")
		viper.SetDefault("http.read_timeout", "10s")
		viper.SetDefault("http.read_header_timeout", "5s")
		viper.SetDefault("http.write_timeout", "30s")
		viper.SetDefault("http.idle_timeout", "60s")
		viper.SetDefault("http.max_header_bytes", 1<<19)
	})
}
