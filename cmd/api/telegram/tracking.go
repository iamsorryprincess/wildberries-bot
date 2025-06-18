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
	deleteTrackingURL     = "/deletetracking/"

	buttonsPerMessage = 20
	buttonsPerRow     = 4
	buttonsRowCount   = 5
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]model.Category, error)
}

type SizeRepository interface {
	GetSizesInfo(ctx context.Context, categoryID uint64) ([]model.SizeInfo, error)
	GetSizeCategoryInfo(ctx context.Context, sizeID uint64, categoryID uint64) (model.SizeCategoryInfo, error)
}

type TrackingRepository interface {
	AddTracking(ctx context.Context, settings model.TrackingSettings) error
	GetTrackingSettingsInfo(ctx context.Context, chatID int64) ([]model.TrackingSettingsInfo, error)
	DeleteTrackingSettings(ctx context.Context, chatID int64, sizeID uint64, categoryID uint64) error
}

type trackingHandler struct {
	logger log.Logger

	categoryRepository CategoryRepository
	sizeRepository     SizeRepository
	trackingRepository TrackingRepository
}

func newTrackingHandler(
	logger log.Logger,
	categoryRepository CategoryRepository,
	sizeRepository SizeRepository,
	trackingRepository TrackingRepository,
) *trackingHandler {
	return &trackingHandler{
		logger:             logger,
		categoryRepository: categoryRepository,
		sizeRepository:     sizeRepository,
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
			CallbackData: fmt.Sprintf("%s%d:%s:%s", trackingCategoriesURL, category.ID, category.Title, category.Emoji),
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
	categoryID, err := strconv.ParseUint(data[:index], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "ShowSizeTrackingOptions").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse category_id from callback query data")
		return
	}

	categoryInfoStr := data[index+1:]
	infoIndex := strings.Index(categoryInfoStr, ":")
	categoryTitle := categoryInfoStr[:infoIndex]
	categoryEmoji := categoryInfoStr[infoIndex+1:]

	sizes, err := h.sizeRepository.GetSizesInfo(ctx, categoryID)
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
					Text:         size.Name,
					CallbackData: fmt.Sprintf("%s%d/%d:%s:%s", showDiffPricesURL, categoryID, size.ID, categoryTitle, categoryEmoji),
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
	categoryID, err := strconv.ParseUint(data[:index], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "ShowDiffPriceOptions").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse category_id from callback query data")
		return
	}

	sizeStr := data[index+1:]
	emojiSeparatorIndex := strings.Index(sizeStr, ":")

	sizeID, err := strconv.ParseUint(sizeStr[:emojiSeparatorIndex], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "ShowDiffPriceOptions").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse size_id from callback query data")
		return
	}

	categoryDataStr := sizeStr[emojiSeparatorIndex+1:]
	dataIndex := strings.Index(categoryDataStr, ":")
	categoryTitle := categoryDataStr[:dataIndex]
	categoryEmoji := categoryDataStr[dataIndex+1:]

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   fmt.Sprintf("Выберите процент снижения цен на %s %s для уведомления:", categoryTitle, categoryEmoji),
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "5%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "5", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "10%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "10", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "15%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "15", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "20%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "20", sizeID, categoryTitle, categoryEmoji, categoryID)},
				},
				{
					{Text: "25%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "25", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "30%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "30", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "35%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "35", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "40%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "40", sizeID, categoryTitle, categoryEmoji, categoryID)},
				},
				{
					{Text: "45%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "45", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "50%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "50", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "55%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "55", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "60%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "60", sizeID, categoryTitle, categoryEmoji, categoryID)},
				},
				{
					{Text: "65%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "65", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "70%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "70", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "75%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "75", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "80%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "80", sizeID, categoryTitle, categoryEmoji, categoryID)},
				},
				{
					{Text: "85%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "85", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "90%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "90", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "95%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "95", sizeID, categoryTitle, categoryEmoji, categoryID)},
					{Text: "100%", CallbackData: fmt.Sprintf("%s%s/%d:%s:%s:%d", addTrackingURL, "100", sizeID, categoryTitle, categoryEmoji, categoryID)},
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

	sizeID, err := strconv.ParseUint(trackingParams[0], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "AddTracking").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse size_id from callback query data")
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

	categoryTitle := trackingParams[1]
	categoryEmoji := trackingParams[2]
	categoryID, err := strconv.ParseUint(trackingParams[3], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "AddTracking").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse category_id from callback query data")
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

	trackingSettings := model.TrackingSettings{
		ChatID:     chatID,
		SizeID:     sizeID,
		CategoryID: categoryID,
		DiffValue:  diffPercent,
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
<b>Снижение цены</b>: <i>%d%%</i> ⬇️`

	sizeData, err := h.sizeRepository.GetSizeCategoryInfo(ctx, sizeID, categoryID)
	if err != nil {
		h.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed get size category info")
		sizeData.Name = "не удалось получить данные :С"
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      fmt.Sprintf(messageText, categoryTitle, categoryEmoji, sizeData.Name, diffPercent),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "AddTrackingSize").
			Int64("chat_id", chatID).
			Msg("failed send message")
	}
}

func (h *trackingHandler) ShowTrackingSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	trackingSettings, err := h.trackingRepository.GetTrackingSettingsInfo(ctx, chatID)
	if err != nil {
		h.logger.Error().Err(err).Str("handler", "ShowTrackingSettings").Msg("get tracking settings failed")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению пока данный функционал недоступен, попробуйте позже :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}

		return
	}

	if len(trackingSettings) == 0 {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "На данный момент у вас отсутствуют настройки отслеживания",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	const messageText = `<b>Категория</b>: <i>%s</i> %s
<b>Размер</b>: <i>%s</i> 📏
<b>Снижение цены</b>: <i>%d%%</i> ⬇️`

	var sb strings.Builder
	sb.WriteString("Ваши текущие настройки отслеживания:")

	for _, settings := range trackingSettings {
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf(messageText, settings.CategoryTitle, settings.CategoryEmoji, settings.Size, settings.DiffPercent))
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      sb.String(),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "ShowTrackingSettings").
			Int64("chat_id", chatID).
			Msg("failed send message")
	}
}

func (h *trackingHandler) ShowDeleteTrackingSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	trackingSettings, err := h.trackingRepository.GetTrackingSettingsInfo(ctx, chatID)
	if err != nil {
		h.logger.Error().Err(err).Str("handler", "ShowDeleteTrackingSettings").Msg("get tracking settings failed")

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению пока данный функционал недоступен, попробуйте позже :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowDeleteTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}

		return
	}

	if len(trackingSettings) == 0 {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "На данный момент у вас отсутствуют настройки отслеживания",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "ShowTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	const msgText = "❌   %s %s %s 📏 %d%% ⬇️"

	var rows [][]models.InlineKeyboardButton
	for _, settings := range trackingSettings {
		rows = append(rows, []models.InlineKeyboardButton{{
			Text:         fmt.Sprintf(msgText, settings.CategoryTitle, settings.CategoryEmoji, settings.Size, settings.DiffPercent),
			CallbackData: fmt.Sprintf("%s%d:%d", deleteTrackingURL, settings.CategoryID, settings.SizeID),
		}})
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Выберите настройку для удаления:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		},
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "ShowDeleteTrackingSettings").
			Int64("chat_id", chatID).
			Msg("failed send message")
	}
}

func (h *trackingHandler) DeleteTrackingSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		h.logger.Error().Str("handler", "DeleteTrackingSettings").Msg("callback query is empty")
		return
	}

	data, isFound := strings.CutPrefix(update.CallbackQuery.Data, deleteTrackingURL)
	if !isFound {
		h.logger.Error().Str("handler", "DeleteTrackingSettings").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't extract data from callback query data")
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	values := strings.Split(data, ":")
	if len(values) < 2 {
		h.logger.Error().Str("handler", "DeleteTrackingSettings").Msg("get tracking settings failed")

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению пока данный функционал недоступен, попробуйте позже :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "DeleteTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}

		return
	}

	categoryID, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "DeleteTrackingSettings").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse category_id from callback query data")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению не удалось удалить настройку отслеживания, попробуйте позже, мы уже чиним поломку :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "DeleteTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	sizeID, err := strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		h.logger.Error().Str("handler", "DeleteTrackingSettings").
			Str("callback_data", update.CallbackQuery.Data).
			Msg("can't parse size_id from callback query data")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению не удалось удалить настройку отслеживания, попробуйте позже, мы уже чиним поломку :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "DeleteTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	if err = h.trackingRepository.DeleteTrackingSettings(ctx, chatID, sizeID, categoryID); err != nil {
		h.logger.Error().Err(err).
			Str("handler", "DeleteTrackingSettings").
			Int64("chat_id", chatID).
			Msg("failed delete tracking settings")
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "К сожалению не удалось удалить настройку отслеживания, попробуйте позже, мы уже чиним поломку :С",
		})
		if err != nil {
			h.logger.Error().Err(err).
				Str("handler", "DeleteTrackingSettings").
				Int64("chat_id", chatID).
				Msg("failed send message")
		}
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Настройка успешно удалена",
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("handler", "DeleteTrackingSettings").
			Int64("chat_id", chatID).
			Msg("failed send message")
	}
}
