package agent

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// waitForShutdown ожидает сигнал завершения и отменяет контекст
func waitForShutdown(cancel context.CancelFunc) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigint
		log.Println("Получен сигнал завершения. Прекращение работы...")
		cancel()
	}()
}
