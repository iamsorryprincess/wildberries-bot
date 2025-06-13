package app

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/config"
	httpapp "github.com/iamsorryprincess/wildberries-bot/cmd/api/http"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/repository"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/service"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

const serviceName = "api"

type App struct {
	config config.Config
	logger log.Logger

	closeStack *background.CloseStack
	appErrors  *background.ErrorsChannel

	ctx context.Context

	mysqlConn *mysql.Connection

	productRepository *repository.MysqlProductRepository

	productClient *httpapp.ProductClient

	productUpdateService *service.ProductUpdateService

	worker *background.Worker

	httpServer *http.Server
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
	a.appErrors = background.NewErrorsChannel()

	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	defer cancel()

	if err = a.initDatabases(); err != nil {
		return
	}

	a.initRepositories()

	a.initServices()

	a.initWorkers()

	a.initHTTP()

	a.logger.Info().Msg("service started")

	s, err := background.Wait(a.appErrors)
	if err != nil {
		a.logger.Error().Err(err).Msg("service failed")
		return
	}

	a.logger.Info().Str("signal", s.String()).Msg("service stopped")
}

func (a *App) initDatabases() error {
	var err error

	a.mysqlConn, err = mysql.NewConnection(a.logger, a.config.MysqlConfig, a.closeStack)
	if err != nil {
		a.logger.Error().Err(err).Msg("mysql connect failed")
		return err
	}

	a.logger.Info().Msg("mysql connected successful")
	return nil
}

func (a *App) initRepositories() {
	a.productRepository = repository.NewMysqlProductRepository(a.mysqlConn)
}

func (a *App) initServices() {
	a.productClient = httpapp.NewProductClient(a.logger, a.config.ProductsClientConfig)
	a.productUpdateService = service.NewUpdateService(a.logger, a.productClient, a.productRepository)
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger, a.closeStack)
	//a.worker.RunWithInterval(a.ctx, "update products", time.Minute*15, a.productUpdateService.Update)
}

func (a *App) initHTTP() {
	handler := httpapp.NewHandler(a.logger)
	a.httpServer = http.NewServer(a.logger, a.config.HTTPConfig, a.closeStack, a.appErrors, handler)
	a.httpServer.Start()
}
