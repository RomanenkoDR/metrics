package agent

import (
	"context"
	"go.uber.org/zap"
	"time"

	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
)

func Run() {
	// Логируем старт приложения
	logger.Info("Начало основного приложения")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов: ", zap.Any("err", err))
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
	logger.Info("Инициализация хранилища успешна. Начало работы")

	// Запускаем основной цикл
	for {
		select {
		case <-pollTicker.C:
			logger.Debug("Сбор метрик")
			ReadMemStats(&memStorage)

		case <-reportTicker.C:
			logger.Debug("Отправка метрик")
			send := Retry(func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
				return ProcessBatch(ctx, serverAddress, cfg.CryptoKey, m)
			}, 3, 1*time.Second)

			err := send(context.Background(), cfg.ServerAddress, memStorage)
			if err != nil {
				logger.DebugLogger.Sugar().Error("Не удалось обработать пакет метрик: ", err)
			}
			logger.Info("Метрики отправлены на сервер")
		}
	}
}
