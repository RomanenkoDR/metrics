package agent

import (
	"context"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"time"
)

// Run запускает основной цикл агента. Включает в себя сбор метрик,
// периодическую отправку на сервер и управление интервалами выполнения операций.
func Run() {
	logger.Info("Начало основного приложения")

	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов: ", zap.Any("err", err))
	}

	logger.Info(fmt.Sprintf("флаг на агенте: ", cfg.Key))

	if cfg.Key != "" {
		Encrypt = true
		Key = []byte(cfg.Key)
	}

	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	memStorage := storage.New()
	logger.Info("Инициализация хранилища успешна. Начало работы")

	for {
		select {
		case <-pollTicker.C:
			logger.Debug("Сбор метрик")
			ReadMemStats(&memStorage)

		case <-reportTicker.C:
			logger.Debug("Отправка метрик")
			send := Retry(ProcessBatch, 3, 1*time.Second)
			err := send(context.Background(), cfg.ServerAddress, memStorage)
			if err != nil {
				logger.DebugLogger.Sugar().Error("Не удалось обработать пакет метрик: ", err)
			}
			logger.Info("Метрики отправлены на сервер")
		}
	}
}
