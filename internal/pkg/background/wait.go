package background

import (
	"os"
	"os/signal"
	"syscall"
)

type AppErrors interface {
	Push(err error)
	Errors() <-chan error
}

func Wait(appErrors AppErrors) (os.Signal, error) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-exit:
		signal.Stop(exit)
		return s, nil
	case err := <-appErrors.Errors():
		return nil, err
	}
}
