package http

import "time"

type ServerConfig struct {
	Port int `config:"port"`

	ShutdownTimeout time.Duration `config:"shutdown_timeout"`

	ReadTimeout       time.Duration `config:"read_timeout"`
	ReadHeaderTimeout time.Duration `config:"read_header_timeout"`
	WriteTimeout      time.Duration `config:"write_timeout"`
	IdleTimeout       time.Duration `config:"idle_timeout"`
	MaxHeaderBytes    int           `config:"max_header_bytes"`
}

type ClientConfig struct {
	Timeout       time.Duration `config:"timeout"`
	DialTimeout   time.Duration `config:"dial_timeout"`
	DialKeepAlive time.Duration `config:"dial_keep_alive"`

	MaxIdleConns        int `config:"max_idle_conns"`
	MaxIdleConnsPerHost int `config:"max_idle_conns_per_host"`

	IdleConnTimeout       time.Duration `config:"idle_conn_timeout"`
	TLSHandshakeTimeout   time.Duration `config:"tls_handshake_timeout"`
	ResponseHeaderTimeout time.Duration `config:"response_header_timeout"`
	ExpectContinueTimeout time.Duration `config:"expect_continue_timeout"`
}
