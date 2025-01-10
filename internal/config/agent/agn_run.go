package agent

import (
	"context"
	"time"

	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
)

func Run() {
	// Логируем старт приложения
	logger.DebugLogger.Sugar().Info("Начало основного приложения")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.DebugLogger.Sugar().Fatal("Ошибка разбора флагов: ", err)
	}

	if cfg.Key != "" {
		Encrypt = true
		Key = []byte(cfg.Key)
	}

	// Создаем тикеры
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	// Инициализируем хранилище метрик
	memStorage := storage.New()
	logger.DebugLogger.Sugar().Info("Инициализация хранилища успешна. Начало работы")

	// Запускаем основной цикл
	for {
		select {
		case <-pollTicker.C:
			logger.DebugLogger.Sugar().Debug("Сбор метрик")
			ReadMemStats(&memStorage)

		case <-reportTicker.C:
			logger.DebugLogger.Sugar().Debug("Отправка метрик")
			send := Retry(ProcessBatch, 3, 1*time.Second)
			err := send(context.Background(), cfg.ServerAddress, memStorage)
			if err != nil {
				logger.DebugLogger.Sugar().Error("Не удалось обработать пакет метрик: ", err)
			}
		}
	}
}
