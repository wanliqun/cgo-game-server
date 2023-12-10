package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// GracefulShutdown supports to clean up goroutines after termination signal captured.
func GracefulShutdown(wg *sync.WaitGroup, shutdown func()) {
	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	// Wait for SIGTERM to be captured
	<-termChan

	// Shutdown to notify active goroutines to clean up.
	shutdown()

	wg.Wait()
}
