package http

import "time"

type Config struct {
	Port int `config:"port"`

	ShutdownTimeout time.Duration `config:"shutdown_timeout"`

	ReadTimeout       time.Duration `config:"read_timeout"`
	ReadHeaderTimeout time.Duration `config:"read_header_timeout"`
	WriteTimeout      time.Duration `config:"write_timeout"`
	IdleTimeout       time.Duration `config:"idle_timeout"`
	MaxHeaderBytes    int           `config:"max_header_bytes"`
}
