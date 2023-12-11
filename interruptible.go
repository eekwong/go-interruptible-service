package interruptible

import (
	"os"
	"os/signal"
	"syscall"
)

type Service interface {
	Run() (Stop, error)
}

type Stop func() error

func Run(s Service) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	stopFunc, err := s.Run()
	if err != nil {
		if stopFunc != nil {
			stopFunc()
		}
		return err
	}
	<-c
	if stopFunc != nil {
		return stopFunc()
	}
	return nil
}
