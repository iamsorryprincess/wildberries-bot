package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

const (
	trackingCategoriesURL = "/trackingcategories/"
	showDiffPricesURL     = "/showdiffprices/"
	addTrackingURL        = "/addtracking/"

	buttonsPerMessage = 20
	buttonsPerRow     = 4
	buttonsRowCount   = 5
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
}

type ProductRepository interface {
	GetSizes(ctx context.Context, category string) ([]string, error)
}

type TrackingRepository interface {
	AddTracking(ctx context.Context, settings model.TrackingSettings) error
}

type trackingHandler struct {
	logger log.Logger

	categoryRepository CategoryRepository
	productRepository  ProductRepository
	trackingRepository TrackingRepository
}

func newTrackingHandler(
	logger log.Logger,
	categoryRepository CategoryRepository,
	productRepository ProductRepository,
	trackingRepository TrackingRepository,
) *trackingHandler {
	return &trackingHandler{
		logger:             logger,
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
		trackingRepository: trackingRepository,
	}
}

func (h *trackingHandler) ShowCategoryTrackingOptions(ctx context.Context, b *bot.Bot, update *models.Update) {
	categories, err := h.categoryRepository.GetCategories(ctx)
	if err != nil || len(categories) == 0 {
		h.logger.Error().Err(err).Str("handler", "ShowCategoryTrackingOptions").Msg("get categories failed")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "К сожалению пока данный функционал недоступен, попробуйте позже :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowCategoryTrackingOptions").
				Int64("chat_id", update.Message.Chat.ID).
				Msg("failed send message")
		}

		return
	}

	keyboardButtons := make([]models.InlineKeyboardButton, len(categories))
	for i, category := range categories {
		keyboardButtons[i] = models.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s %s", category.Emoji, category.Title),
			CallbackData: fmt.Sprintf("%s%s:%s:%s", trackingCategoriesURL, category.Name, category.Title, category.Emoji),
		}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Выберите категорию товара для отслеживания:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				keyboardButtons,
			},
		},
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "ShowCategoryTrackingOptions").
			Int64("chat_id", update.Message.Chat.ID).
			Msg("failed send message")
	}
}

func (h *trackingHandler) ShowSizeTrackingOptions(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		h.logger.Error().Str("handler", "ShowSizeTrackingOptions").Msg("callback query is empty")
		return
	}

	data, isFound := strings.CutPrefix(update.CallbackQuery.Data, trackingCategoriesURL)
	if !isFound {
		h.logger.Error().Str("handler", "ShowSizeTrackingOptions").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't extract data from callback query data")
		return
	}

	index := strings.Index(data, ":")
	category := data[:index]
	categoryInfoStr := data[index+1:]
	infoIndex := strings.Index(categoryInfoStr, ":")
	categoryTitle := categoryInfoStr[:infoIndex]
	categoryEmoji := categoryInfoStr[infoIndex+1:]

	sizes, err := h.productRepository.GetSizes(ctx, category)
	if err != nil || len(sizes) == 0 {
		h.logger.Error().Err(err).Str("handler", "ShowSizeTrackingOptions").Msg("get sizes failed")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "К сожалению для данной категории пока нет информации о товарах, попробуйте позже :)",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowSizeTrackingOptions").
				Int64("chat_id", update.CallbackQuery.Message.Message.Chat.ID).
				Msg("failed send message")
		}

		return
	}

	itemIndex := 0
	for itemIndex < len(sizes) {
		isBreak := false
		var rows [][]models.InlineKeyboardButton

		for rowIndex := 0; rowIndex < buttonsRowCount; rowIndex++ {
			if isBreak {
				break
			}

			var row []models.InlineKeyboardButton

			for columnIndex := 0; columnIndex < buttonsPerRow; columnIndex++ {
				size := sizes[itemIndex]

				row = append(row, models.InlineKeyboardButton{
					Text:         size,
					CallbackData: fmt.Sprintf("%s%s/%s:%s:%s", showDiffPricesURL, category, size, categoryTitle, categoryEmoji),
				})

				itemIndex++

				if itemIndex == len(sizes) {
					isBreak = true
					break
				}
			}

			rows = append(rows, row)
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   fmt.Sprintf("Выберите размер для категории %s %s:", categoryTitle, categoryEmoji),
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: rows,
			},
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowSizeTrackingOptions").
				Int64("chat_id", update.CallbackQuery.Message.Message.Chat.ID).
				Msg("failed send message")
		}
	}
}

func (h *trackingHandler) ShowDiffPriceOptions(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		h.logger.Error().Str("handler", "ShowDiffPriceOptions").Msg("callback query is empty")
		return
	}

	data, isFound := strings.CutPrefix(update.CallbackQuery.Data, showDiffPricesURL)
	if !isFound {
		h.logger.Error().Str("handler", "ShowDiffPriceOptions").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't extract data from callback query data")
		return
	}

	index := strings.Index(data, "/")
	category := data[:index]
	sizeStr := data[index+1:]
	emojiSeparatorIndex := strings.Index(sizeStr, ":")
	size := sizeStr[:emojiSeparatorIndex]
	categoryDataStr := sizeStr[emojiSeparatorIndex+1:]
	dataIndex := strings.Index(categoryDataStr, ":")
	categoryTitle := categoryDataStr[:dataIndex]
	categoryEmoji := categoryDataStr[dataIndex+1:]

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   fmt.Sprintf("Выберите процент снижения цен на %s %s для уведомления:", categoryTitle, categoryEmoji),
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "5%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "5", size, categoryTitle, categoryEmoji, category)},
					{Text: "10%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "10", size, categoryTitle, categoryEmoji, category)},
					{Text: "15%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "15", size, categoryTitle, categoryEmoji, category)},
					{Text: "20%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "20", size, categoryTitle, categoryEmoji, category)},
				},
				{
					{Text: "25%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "25", size, categoryTitle, categoryEmoji, category)},
					{Text: "30%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "30", size, categoryTitle, categoryEmoji, category)},
					{Text: "35%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "35", size, categoryTitle, categoryEmoji, category)},
					{Text: "40%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "40", size, categoryTitle, categoryEmoji, category)},
				},
				{
					{Text: "45%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "45", size, categoryTitle, categoryEmoji, category)},
					{Text: "50%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "50", size, categoryTitle, categoryEmoji, category)},
					{Text: "55%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "55", size, categoryTitle, categoryEmoji, category)},
					{Text: "60%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "60", size, categoryTitle, categoryEmoji, category)},
				},
				{
					{Text: "65%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "65", size, categoryTitle, categoryEmoji, category)},
					{Text: "70%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "70", size, categoryTitle, categoryEmoji, category)},
					{Text: "75%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "75", size, categoryTitle, categoryEmoji, category)},
					{Text: "80%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "80", size, categoryTitle, categoryEmoji, category)},
				},
				{
					{Text: "85%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "85", size, categoryTitle, categoryEmoji, category)},
					{Text: "90%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "90", size, categoryTitle, categoryEmoji, category)},
					{Text: "95%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "95", size, categoryTitle, categoryEmoji, category)},
					{Text: "100%", CallbackData: fmt.Sprintf("%s%s/%s:%s:%s:%s", addTrackingURL, "100", size, categoryTitle, categoryEmoji, category)},
				},
			},
		},
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "ShowDiffPriceOptions").
			Int64("chat_id", update.CallbackQuery.Message.Message.Chat.ID).
			Msg("failed send message")
	}
}

func (h *trackingHandler) AddTracking(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		h.logger.Error().Str("handler", "AddTrackingSize").Msg("callback query is empty")
		return
	}

	data, isFound := strings.CutPrefix(update.CallbackQuery.Data, addTrackingURL)
	if !isFound {
		h.logger.Error().Str("handler", "AddTrackingSize").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't extract data from callback query data")
		return
	}

	index := strings.Index(data, "/")
	diffPercentStr := data[:index]
	diffPercent, err := strconv.Atoi(diffPercentStr)
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "AddTrackingSize").
			Str("value", diffPercentStr).
			Msg("failed parse diff price value")
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	sizeStr := data[index+1:]
	trackingParams := strings.Split(sizeStr, ":")
	if len(trackingParams) < 4 {
		h.logger.Error().Str("params_str", sizeStr).Msg("invalid tracking params count")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "К сожалению не удалось добавить настройку отслеживания, попробуйте позже, мы уже чиним поломку :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "AddTrackingSize").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	size := trackingParams[0]
	categoryTitle := trackingParams[1]
	categoryEmoji := trackingParams[2]
	categoryName := trackingParams[3]

	trackingSettings := model.TrackingSettings{
		ChatID:    chatID,
		Size:      size,
		Category:  categoryName,
		DiffValue: diffPercent,
	}

	if err = h.trackingRepository.AddTracking(ctx, trackingSettings); err != nil {
		h.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed add tracking settings")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "К сожалению не удалось добавить настройку отслеживания, попробуйте позже, мы уже чиним поломку :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "AddTrackingSize").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	const messageText = `Вы добавили настройки отслеживания для следующих параметров:
<b>Категория</b>: <i>%s</i> %s
<b>Размер</b>: <i>%s</i> 📏
<b>Снижение цены</b>: <i>%d%%</i>`

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      fmt.Sprintf(messageText, categoryTitle, categoryEmoji, size, diffPercent),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "AddTrackingSize").
			Int64("chat_id", chatID).
			Msg("failed send message")
	}
}
