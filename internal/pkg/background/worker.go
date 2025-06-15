package background

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type HandlerFunc func(context.Context) error

type Worker struct {
	logger log.Logger

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewWorker(logger log.Logger, closerStack CloserStack) *Worker {
	w := &Worker{
		logger: logger,
	}

	closerStack.Push(w)
	return w
}

func (w *Worker) Run(ctx context.Context, name string, handler HandlerFunc) {
	w.mu.Lock()
	w.wg.Add(1)
	w.mu.Unlock()

	go func(ctx context.Context, handler HandlerFunc) {
		defer w.wg.Done()
		start := time.Now()

		if err := handler(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
				w.logger.Info().Str("worker", name).Msg("worker canceled")
				return
			}
			w.logger.Error().Err(err).Str("worker", name).Msg("worker error")
		}

		w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
	}(ctx, handler)
}

func (w *Worker) RunWithInterval(ctx context.Context, name string, interval time.Duration, handler HandlerFunc) {
	w.mu.Lock()
	w.wg.Add(1)
	w.mu.Unlock()

	go func(ctx context.Context, interval time.Duration, handler HandlerFunc) {
		defer w.wg.Done()

		start := time.Now()
		if err := handler(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
				w.logger.Info().Str("worker", name).Msg("worker canceled")
				return
			}
			w.logger.Error().Err(err).Str("worker", name).Msg("worker error")
		}
		w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				w.logger.Info().Str("worker", name).Msg("worker stopped")
				return
			case <-timer.C:
				start = time.Now()
				if err := handler(ctx); err != nil {
					if errors.Is(err, context.Canceled) {
						w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
						w.logger.Info().Str("worker", name).Msg("worker canceled")
						return
					}
					w.logger.Error().Err(err).Str("worker", name).Msg("worker error")
				}
				w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
				timer.Reset(interval)
			}
		}
	}(ctx, interval, handler)
}

func (w *Worker) Close() error {
	w.wg.Wait()
	return nil
}
