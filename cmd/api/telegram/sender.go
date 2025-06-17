package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/telegram"
)

type Sender struct {
	client *telegram.BotClient
}

func NewSender(client *telegram.BotClient) *Sender {
	return &Sender{
		client: client,
	}
}

func (s *Sender) Send(ctx context.Context, message model.TrackingResult) error {
	const messageText = `<b>%s</b>

<a href="%s">Ссылка на товар</a>

<b>Размер:</b> %s

<b>Старая цена:</b> %s

<b>Новая цена:</b> %s

<b>Снижение цены:</b> %d%%`
	_, err := s.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: message.ChatID,
		Text: fmt.Sprintf(
			messageText,
			message.ProductName,
			message.ProductURL,
			message.Size,
			message.PreviousPrice,
			message.CurrentPrice,
			message.DiffPercent,
		),
		ParseMode: models.ParseModeHTML,
	})
	return err
}
