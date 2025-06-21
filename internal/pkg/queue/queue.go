package queue

import "context"

type Handler[TMessage any] func(ctx context.Context, messages []TMessage) error

type Queue[TMessage any] interface {
	Push(ctx context.Context, message TMessage) error
}
