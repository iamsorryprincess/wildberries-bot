package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type startHandler struct {
	logger log.Logger
}

func NewStartHandlerOption(logger log.Logger) bot.Option {
	handler := &startHandler{
		logger: logger,
	}
	return bot.WithDefaultHandler(handler.Handle)
}

func (h *startHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	const message = `Я могу помочь вам создать и управлять настройками отслеживания цен товаров Wildberries.

Вы можете управлять мной, отправляя следующие команды:

/addtracking - добавляет отслеживание
/edittracking - редактирует настройки отслеживания
/deletetracking - удаляет отслеживание
/showtracking - показывает текущие настройки отслеживания`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Int64("chat_id", update.Message.Chat.ID).
			Str("handler", "start").
			Msg("failed telegram send message")
	}
}
