package background

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type HandlerFunc func(context.Context) error

type Worker struct {
	logger log.Logger
	wg     sync.WaitGroup
}

func NewWorker(logger log.Logger) *Worker {
	return &Worker{
		logger: logger,
	}
}

func (w *Worker) Run(ctx context.Context, name string, handler HandlerFunc) {
	w.wg.Add(1)

	go func(ctx context.Context, handler HandlerFunc) {
		defer w.wg.Done()
		start := time.Now()
		defer w.recover(name)

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
	w.wg.Add(1)

	go func(ctx context.Context, interval time.Duration, handler HandlerFunc) {
		defer w.wg.Done()
		timer := time.NewTimer(interval)
		defer timer.Stop()

		if err := w.handle(ctx, name, handler); err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				w.logger.Info().Str("worker", name).Msg("worker stopped")
				return
			case <-timer.C:
				if err := w.handle(ctx, name, handler); err != nil {
					return
				}
				timer.Reset(interval)
			}
		}
	}(ctx, interval, handler)
}

func (w *Worker) Close() error {
	w.wg.Wait()
	return nil
}

func (w *Worker) handle(ctx context.Context, name string, handler HandlerFunc) error {
	defer w.recover(name)
	start := time.Now()
	if err := handler(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
			w.logger.Info().Str("worker", name).Msg("worker canceled")
			return err
		}
		w.logger.Error().Err(err).Str("worker", name).Msg("worker error")
	}
	w.logger.Debug().Str("worker", name).Str("duration", time.Since(start).String()).Send()
	return nil
}

func (w *Worker) recover(name string) {
	if rvr := recover(); rvr != nil {
		event := w.logger.Error().
			Str("worker", name).
			Str("stack", string(debug.Stack()))

		err, ok := rvr.(error)
		if ok {
			event = event.Err(err)
		} else {
			event = event.Interface("panic", rvr)
		}

		event.Msg("recovered from panic")
	}
}
