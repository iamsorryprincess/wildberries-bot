package http

import (
	"net"
	"net/http"
)

func NewClient(config ClientConfig) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   config.DialTimeout,
				KeepAlive: config.DialKeepAlive,
			}).DialContext,

			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,

			IdleConnTimeout:       config.Timeout,
			TLSHandshakeTimeout:   config.Timeout,
			ResponseHeaderTimeout: config.Timeout,
			ExpectContinueTimeout: config.Timeout,
		},

		Timeout: config.Timeout,
	}
}
