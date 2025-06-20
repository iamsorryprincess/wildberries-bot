package config

import (
	"time"

	httpapp "github.com/iamsorryprincess/wildberries-bot/cmd/api/http"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/config"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/telegram"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `config:"loglevel"`

	ParseInterval time.Duration `config:"parse_interval"`

	MysqlConfig mysql.Config `config:"mysql"`

	HTTPClientConfig http.ClientConfig `config:"http_client"`

	ProductsClientConfig httpapp.ProductClientConfig `config:"products_client"`

	TelegramConfig telegram.Config `config:"telegram"`

	HTTPConfig http.ServerConfig `config:"http"`
}

func Init() (Config, error) {
	return config.Load[Config](func() {
		viper.SetDefault("loglevel", "info")
		viper.SetDefault("parse_interval", "15m")

		viper.SetDefault("mysql.max_open_connections", 5)
		viper.SetDefault("mysql.max_idle_connections", 5)
		viper.SetDefault("mysql.connection_max_lifetime", 5*time.Minute)
		viper.SetDefault("mysql.connection_max_idle_time", 5*time.Minute)

		viper.SetDefault("http_client.timeout", 30*time.Second)
		viper.SetDefault("http_client.dial_timeout", 5*time.Second)
		viper.SetDefault("http_client.dial_keep_alive", 30*time.Second)
		viper.SetDefault("http_client.max_idle_conns", 10)
		viper.SetDefault("http_client.max_idle_conns_per_host", 10)
		viper.SetDefault("http_client.idle_conn_timeout", 90*time.Second)
		viper.SetDefault("http_client.tls_handshake_timeout", 10*time.Second)
		viper.SetDefault("http_client.response_header_timeout", 5*time.Second)
		viper.SetDefault("http_client.expect_continue_timeout", 1*time.Second)

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
