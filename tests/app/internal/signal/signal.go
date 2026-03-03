// Package signal contains helpers for graceful shutdown on OS signals.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler waits for SIGINT or SIGTERM and invokes the provided cancel function.
func Handler(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	doneChan := make(chan bool, 1)

	defer cancel()

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		doneChan <- true
	}()

	<-doneChan
}
