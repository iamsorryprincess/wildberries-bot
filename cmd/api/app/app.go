package app

import (
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/config"
	configutils "github.com/iamsorryprincess/wildberries-bot/internal/pkg/config"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

const serviceName = "api"

type App struct {
	config config.Config
	logger log.Logger
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	var err error

	if a.config, err = configutils.Load[config.Config](); err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("failed to load config")
		return
	}

	a.logger = log.New(a.config.LogLevel, serviceName)
}
