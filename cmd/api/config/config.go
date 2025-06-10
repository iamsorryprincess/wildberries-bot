package config

import (
	"time"

	httpapp "github.com/iamsorryprincess/wildberries-bot/cmd/api/http"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/config"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `config:"loglevel"`

	ProductsClientConfig httpapp.ProductClientConfig `config:"products_client"`

	HTTPConfig http.Config `config:"http"`
}

func Init() (Config, error) {
	return config.Load[Config](func() {
		viper.SetDefault("loglevel", "info")

		viper.SetDefault("products_client.retry_count", 3)
		viper.SetDefault("products_client.retry_delay", time.Second)

		viper.SetDefault("http.port", "8080")
		viper.SetDefault("http.read_timeout", "10s")
		viper.SetDefault("http.read_header_timeout", "5s")
		viper.SetDefault("http.write_timeout", "30s")
		viper.SetDefault("http.idle_timeout", "60s")
		viper.SetDefault("http.max_header_bytes", 1<<19)
	})
}
