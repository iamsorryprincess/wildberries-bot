package mysql

import "time"

type Config struct {
	ConnectionString string `config:"connection_string"`

	MaxOpenConnections int `config:"max_open_connections"`
	MaxIdleConnections int `config:"max_idle_connections"`

	ConnectionMaxLifetime time.Duration `config:"connection_max_lifetime"`
	ConnectionMaxIdleTime time.Duration `config:"connection_max_idle_time"`
}
