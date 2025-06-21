package memory

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/queue"
)

type Queue[TMessage any] struct {
	logger log.Logger
	config Config

	wg      sync.WaitGroup
	handler queue.Handler[TMessage]

	exit     chan struct{}
	messages chan TMessage
}

func NewQueue[TMessage any](ctx context.Context, logger log.Logger, config Config, handler queue.Handler[TMessage]) *Queue[TMessage] {
	queue := &Queue[TMessage]{
		logger:   logger,
		config:   config,
		handler:  handler,
		exit:     make(chan struct{}),
		messages: make(chan TMessage, config.BufferSize),
	}

	queue.wg.Add(1)

	go func(ctx context.Context, queue *Queue[TMessage]) {
		defer queue.wg.Done()
		batch := make([]TMessage, 0, queue.config.BatchSize)
		timer := time.NewTimer(queue.config.FlushInterval)

		for {
			select {
			case _, ok := <-queue.exit:
				if !ok {
					queue.logger.Debug().Msg("memory queue stopped")

					if len(batch) > 0 {
						queue.handle(ctx, batch)
					}

					if len(queue.messages) > 0 {
						queue.drain(ctx)
					}

					return
				}
			case msg := <-queue.messages:
				batch = append(batch, msg)

				if len(batch) >= queue.config.BatchSize {
					queue.handle(ctx, batch)
					batch = batch[:0]
				}
			case <-timer.C:
				if len(batch) > 0 {
					queue.handle(ctx, batch)
					batch = batch[:0]
				}
				timer.Reset(queue.config.FlushInterval)
			}
		}
	}(ctx, queue)

	return queue
}

func (q *Queue[TMessage]) Push(_ context.Context, message TMessage) error {
	q.messages <- message
	return nil
}

func (q *Queue[TMessage]) Close() error {
	close(q.exit)
	q.wg.Wait()
	return nil
}

func (q *Queue[TMessage]) handle(ctx context.Context, messages []TMessage) {
	defer q.recover()
	if err := q.handler(ctx, messages); err != nil {
		if errors.Is(err, context.Canceled) {
			q.logger.Info().Msg("memory queue handler canceled")
			return
		}
		q.logger.Error().Err(err).Msg("memory queue failed handle messages")
	}
}

func (q *Queue[TMessage]) drain(ctx context.Context) {
	batch := make([]TMessage, 0, len(q.messages))

loop:
	for {
		select {
		case msg := <-q.messages:
			batch = append(batch, msg)
		default:
			break loop
		}
	}

	q.handle(ctx, batch)
}

func (q *Queue[TMessage]) recover() {
	if rvr := recover(); rvr != nil {
		event := q.logger.Error().Str("stack", string(debug.Stack()))

		err, ok := rvr.(error)
		if ok {
			event = event.Err(err)
		} else {
			event = event.Interface("panic", rvr)
		}

		event.Msg("memory queue handler recovered from panic")
	}
}
