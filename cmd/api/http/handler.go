package http

import (
	"net/http"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http/middleware"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

func NewHandler(logger log.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return middleware.WithHandler(mux, middleware.Recovery(logger))
}
