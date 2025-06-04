package app

import (
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/config"
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
	cfg, err := config.Init()
	if err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("init config failed")
		return
	}

	a.config = cfg
	a.logger = log.New(a.config.LogLevel, serviceName)
}
