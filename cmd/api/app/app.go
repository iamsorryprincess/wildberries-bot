package app

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/config"
	httptransport "github.com/iamsorryprincess/wildberries-bot/cmd/api/http"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/repository"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/service"
	telegramtransport "github.com/iamsorryprincess/wildberries-bot/cmd/api/telegram"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/http"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/telegram"
)

const serviceName = "api"

type App struct {
	config config.Config
	logger log.Logger

	closeStack *background.CloseStack
	appErrors  *background.ErrorsChannel

	ctx context.Context

	mysqlConn *mysql.Connection

	categoryRepository *repository.MysqlCategoryRepository
	productRepository  *repository.MysqlProductRepository
	sizeRepository     *repository.MysqlSizeRepository
	trackingRepository *repository.MysqlTrackingRepository

	productClient *httptransport.ProductClient
	botClient     *telegram.BotClient

	sender *telegramtransport.Sender

	trackingService *service.TrackingService
	productService  *service.ProductService

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
	a.appErrors = background.NewErrorsChannel()

	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	defer cancel()

	if err = a.initDatabases(); err != nil {
		return
	}

	a.initRepositories()

	if err = a.initTelegram(); err != nil {
		return
	}

	a.initWorkers()

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
	a.categoryRepository = repository.NewMysqlCategoryRepository(a.logger, a.mysqlConn)
	a.productRepository = repository.NewMysqlProductRepository(a.logger, a.mysqlConn)
	a.sizeRepository = repository.NewMysqlSizeRepository(a.mysqlConn)
	a.trackingRepository = repository.NewMysqlTrackingRepository(a.logger, a.mysqlConn)
}

func (a *App) initTelegram() error {
	defaultHandlerOption := telegramtransport.NewStartHandlerOption(a.logger, a.trackingRepository)

	botClient, err := telegram.NewBotClient(a.config.TelegramConfig, defaultHandlerOption)
	if err != nil {
		a.logger.Error().Err(err).Msg("telegram bot client init failed")
		return err
	}

	a.botClient = botClient
	telegramtransport.InitHandlers(a.logger, a.botClient, a.categoryRepository, a.sizeRepository, a.trackingRepository)
	a.sender = telegramtransport.NewSender(a.botClient)
	a.botClient.Start(a.ctx, a.closeStack)
	return nil
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger, a.closeStack)
	a.productClient = httptransport.NewProductClient(a.logger, a.config.ProductsClientConfig, http.NewClient(a.config.HTTPClientConfig))
	a.trackingService = service.NewTrackingService(a.logger, a.trackingRepository, a.sender)

	a.productService = service.NewProductService(
		a.logger,
		a.productClient,
		a.categoryRepository,
		a.productRepository,
		a.trackingService,
	)

	a.worker.RunWithInterval(a.ctx, "run updates", a.config.ParseInterval, a.productService.RunUpdateWorkers)
}
