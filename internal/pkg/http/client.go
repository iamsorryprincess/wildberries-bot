package http

import (
	"net"
	"net/http"
)

func NewClient(config ClientConfig) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   config.DialTimeout,   // Тайм-аут для установления соединения
				KeepAlive: config.DialKeepAlive, // Поддержание соединений активными
			}).DialContext,

			MaxIdleConns:        config.MaxIdleConns,        // Максимальное количество бездействующих соединений в пуле
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost, // Максимальное количество бездействующих соединений на хост

			IdleConnTimeout:       config.Timeout, // Тайм-аут для бездействующих соединений
			TLSHandshakeTimeout:   config.Timeout, // Тайм-аут для TLS рукопожатия
			ResponseHeaderTimeout: config.Timeout, // Тайм-аут ожидания заголовков ответа
			ExpectContinueTimeout: config.Timeout, // Тайм-аут для ожидания продолжения передачи
		},

		Timeout: config.Timeout, // общий таймаут на запрос
	}
}
