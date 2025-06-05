package app

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/config"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

const serviceName = "api"

type App struct {
	config config.Config
	logger log.Logger

	closeStack *background.CloseStack

	ctx context.Context

	worker *background.Worker
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	cfg, err := config.Init()
	if err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("init config failed")
		return
	}

	a.config = cfg
	a.logger = log.New(a.config.LogLevel, serviceName)

	a.closeStack = background.NewCloseStack(a.logger)
	defer a.closeStack.Close()

	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	defer cancel()

	a.initWorkers()

	a.logger.Info().Msg("service started")

	s := background.Wait()

	a.logger.Info().Str("signal", s.String()).Msg("service stopped")
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger, a.closeStack)
}
