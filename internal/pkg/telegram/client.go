package telegram

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/background"
)

type BotClient struct {
	config Config

	*bot.Bot

	wg sync.WaitGroup
}

func NewBotClient(config Config, options ...bot.Option) (*BotClient, error) {
	b, err := bot.New(config.Token, options...)
	if err != nil {
		return nil, err
	}

	botClient := &BotClient{
		config: config,
		Bot:    b,
	}

	return botClient, nil
}

func (c *BotClient) Start(ctx context.Context, closerStack background.CloserStack) {
	closerStack.Push(c)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.Bot.Start(ctx)
	}()
}

func (c *BotClient) Close() error {
	c.wg.Wait()
	return nil
}
