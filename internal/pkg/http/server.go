package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type Server struct {
	logger    log.Logger
	config    Config
	appErrors background.AppErrors
	server    *http.Server
}

func NewServer(logger log.Logger, config Config, closerStack background.CloserStack, appErrors background.AppErrors, handler http.Handler) *Server {
	server := &Server{
		logger:    logger,
		config:    config,
		appErrors: appErrors,
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

	closerStack.Push(server)

	return server
}

func (s *Server) Start() {
	go func() {
		s.logger.Debug().Msgf("starting http server at port %d", s.config.Port)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.appErrors.Push(err)
		}
	}()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
