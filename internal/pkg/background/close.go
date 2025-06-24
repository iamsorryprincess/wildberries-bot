package background

import (
	"io"
	"reflect"
	"sync"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type CloserStack struct {
	logger log.Logger

	mu      sync.Mutex
	closers []io.Closer
}

func NewCloserStack(logger log.Logger) *CloserStack {
	return &CloserStack{
		logger: logger,
	}
}

func (s *CloserStack) Push(closer io.Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closers = append(s.closers, closer)
}

func (s *CloserStack) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := len(s.closers) - 1; i >= 0; i-- {
		closer := s.closers[i]
		t := reflect.TypeOf(closer)

		if err := closer.Close(); err != nil {
			s.logger.Error().Err(err).Str("type", t.String()).Msg("close failed")
			continue
		}

		s.logger.Debug().Str("type", t.String()).Msg("close succeeded")
	}

	s.closers = s.closers[:0]
}
