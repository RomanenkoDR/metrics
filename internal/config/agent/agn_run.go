package agent

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
)

// Run запускает основной процесс агента, включая сбор и отправку метрик.
func Run() {
	logger.Info("Запуск агента")

	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка конфигурации: ", zap.Error(err))
	}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	memStorage := storage.New()
	logger.Info("Агент успешно запущен")

	for {
		select {
		case <-pollTicker.C:
			logger.Debug("Сбор метрик")
			ReadMemStats(&memStorage)
		case <-reportTicker.C:
			logger.Debug("Отправка метрик")
			send := Retry(ProcessBatch, 3, time.Second)
			if err := send(context.Background(), cfg.ServerAddress, memStorage); err != nil {
				logger.Error("Ошибка отправки метрик", zap.Error(err))
			}
		}
	}
}
