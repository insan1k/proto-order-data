package util

import (
	"os"
	"os/signal"
	"syscall"
)

func Signal(q func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-signalChan
		switch s {
		case os.Interrupt, syscall.SIGTERM:
			q()
		}
	}
}

