package background

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait() os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	s := <-exit
	signal.Stop(exit)
	return s
}
