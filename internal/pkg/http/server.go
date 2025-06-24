package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type Server struct {
	logger log.Logger
	config ServerConfig

	fatalErrors chan<- error

	server *http.Server
}

func NewServer(logger log.Logger, config ServerConfig, fatalErrors chan<- error, handler http.Handler) *Server {
	return &Server{
		logger:      logger,
		config:      config,
		fatalErrors: fatalErrors,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", config.Port),
			Handler:           handler,
			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
		},
	}
}

func (s *Server) Start() {
	go func() {
		s.logger.Debug().Msgf("starting http server at port %d", s.config.Port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.fatalErrors <- err
		}
	}()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
