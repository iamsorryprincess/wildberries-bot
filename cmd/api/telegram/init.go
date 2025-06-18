package telegram

import (
	"github.com/go-telegram/bot"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/telegram"
)

func InitHandlers(
	logger log.Logger,
	client *telegram.BotClient,
	categoryRepository CategoryRepository,
	sizeRepository SizeRepository,
	trackingRepository TrackingRepository,
) {
	tracking := newTrackingHandler(logger, categoryRepository, sizeRepository, trackingRepository)

	client.RegisterHandler(bot.HandlerTypeMessageText, "/addtracking", bot.MatchTypeExact, tracking.ShowCategoryTrackingOptions)
	client.RegisterHandler(bot.HandlerTypeCallbackQueryData, trackingCategoriesURL, bot.MatchTypePrefix, tracking.ShowSizeTrackingOptions)
	client.RegisterHandler(bot.HandlerTypeCallbackQueryData, showDiffPricesURL, bot.MatchTypePrefix, tracking.ShowDiffPriceOptions)
	client.RegisterHandler(bot.HandlerTypeCallbackQueryData, addTrackingURL, bot.MatchTypePrefix, tracking.AddTracking)

	client.RegisterHandler(bot.HandlerTypeMessageText, "/showtracking", bot.MatchTypeExact, tracking.ShowTrackingSettings)

	client.RegisterHandler(bot.HandlerTypeMessageText, "/deletetracking", bot.MatchTypeExact, tracking.ShowDeleteTrackingSettings)
	client.RegisterHandler(bot.HandlerTypeCallbackQueryData, deleteTrackingURL, bot.MatchTypePrefix, tracking.DeleteTrackingSettings)
}
