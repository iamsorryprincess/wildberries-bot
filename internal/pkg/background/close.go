package background

import (
	"io"
	"sync"

	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type CloserList interface {
	Add(closer io.Closer)
}

type CloseStack struct {
	logger log.Logger

	mu    sync.Mutex
	elems []io.Closer
}

func NewCloseStack(logger log.Logger) *CloseStack {
	return &CloseStack{
		logger: logger,
	}
}

func (s *CloseStack) Add(closer io.Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elems = append(s.elems, closer)
}

func (s *CloseStack) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := len(s.elems) - 1; i >= 0; i-- {
		closer := s.elems[i]
		if err := closer.Close(); err != nil {
			s.logger.Error().Err(err).Msg("close failed")
		}
		s.logger.Debug().Msg("close succeeded")
	}

	s.elems = s.elems[:0]
}
