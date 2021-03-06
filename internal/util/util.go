package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

//Signal watches for os signal and calls function, preferably quit function
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

//Printer things that come in []byte
func Printer(in chan []byte) (q func()) {
	quit := make(chan struct{})
	q = func() {
		close(quit)
	}
	go func() {
		for {
			select {
			case p := <-in:
				fmt.Printf("%s\n", p)
			case <-quit:
				return
			}
		}
	}()
	return
}
