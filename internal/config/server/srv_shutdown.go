package server

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// setupGracefulShutdown настраивает корректное завершение работы сервера.
func setupShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server, store storage.StorageWriter, metrics *storage.MemStorage) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigint
		logger.DebugLogger.Info("Получен сигнал завершения")
		cancel()

		// Сохраняем текущие метрики перед завершением работы
		if err := store.Write(*metrics); err != nil {
			logger.DebugLogger.Error("Ошибка сохранения данных при завершении", zap.Error(err))
		}

		// Завершаем работу HTTP-сервера
		if err := server.Shutdown(ctx); err != nil {
			logger.DebugLogger.Error("Ошибка завершения HTTP сервера", zap.Error(err))
		}
	}()
}
