package common

import (
	"os"
	"os/signal"
	"syscall"
)

// function to run a cleanup function on signal interruptions such as SIGINT (Ctl+C).
func CloseOnSignalInterrupt(cleanupFunc func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cleanupFunc()
		os.Exit(0)
	}()
}
