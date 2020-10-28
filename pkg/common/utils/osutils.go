package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitSignal block until os.Interrupt, os.Kill, syscall.SIGTERM
func WaitSignal(fn func()) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		stop := false
		select {
		case <-sigChannel:
			stop = true
		}
		if stop {
			fn()
			break
		}
	}
}
