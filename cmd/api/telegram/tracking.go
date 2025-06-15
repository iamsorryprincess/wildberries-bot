package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

const (
	trackingCategoriesURL = "/trackingcategories/"
	trackingSizesURL      = "/trackingsizes/"

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

type trackingHandler struct {
	logger log.Logger

	categoryRepository CategoryRepository
	productRepository  ProductRepository
}

func newTrackingHandler(logger log.Logger, categoryRepository CategoryRepository, productRepository ProductRepository) *trackingHandler {
	return &trackingHandler{
		logger:             logger,
		categoryRepository: categoryRepository,
		productRepository:  productRepository,
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
					CallbackData: fmt.Sprintf("%s%s/%s:%s:%s", trackingSizesURL, category, size, categoryTitle, categoryEmoji),
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

func (h *trackingHandler) AddTrackingSize(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		h.logger.Error().Str("handler", "AddTrackingSize").Msg("callback query is empty")
		return
	}

	data, isFound := strings.CutPrefix(update.CallbackQuery.Data, trackingSizesURL)
	if !isFound {
		h.logger.Error().Str("handler", "AddTrackingSize").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't extract data from callback query data")
		return
	}

	index := strings.Index(data, "/")
	sizeStr := data[index+1:]
	emojiSeparatorIndex := strings.Index(sizeStr, ":")
	size := sizeStr[:emojiSeparatorIndex]
	categoryDataStr := sizeStr[emojiSeparatorIndex+1:]
	dataIndex := strings.Index(categoryDataStr, ":")
	categoryTitle := categoryDataStr[:dataIndex]
	categoryEmoji := categoryDataStr[dataIndex+1:]

	const messageText = `Вы добавили настройки отслеживания для следующих параметров:
<b>Категория</b>: <i>%s</i> %s
<b>Размер</b>: <i>%s</i> 📏`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		Text:      fmt.Sprintf(messageText, categoryTitle, categoryEmoji, size),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "AddTrackingSize").
			Int64("chat_id", update.CallbackQuery.Message.Message.Chat.ID).
			Msg("failed send message")
	}
}
