package main

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/agent"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

func main() {
	// Логируем старт приложения.
	logger.DebugLogger.Sugar().Info("Начало основного приложения")

	// Парсим параметры командной строки с помощью функции из пакета agent.
	cfg, err := agent.ParseOptions()
	if err != nil {
		// Если произошла ошибка при парсинге, логируем фатальную ошибку и завершаем программу.
		logger.DebugLogger.Sugar().Fatal("Ошибка разбора флагов: ", err)
	}

	if cfg.Key != "" {
		agent.Encrypt = true
		agent.Key = []byte(cfg.Key)
	}

	// Создаём тикеры для опроса с интервалами, указанными в конфигурации.
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()

	// Создаём тикеры для отправки с интервалами, указанными в конфигурации.
	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	// Инициализируем новое хранилище данных для метрик.
	m := storage.New()

	// Логируем успешную инициализацию и начало основного цикла программы.
	logger.DebugLogger.Sugar().Info("Инициализация хранилища успешна. Начало основной функции")

	// Основной цикл программы, который работает вечно (пока не завершится).
	for {
		select {

		case <-pollTicker.C:
			logger.DebugLogger.Sugar().Debug("Вызываем сбор метрик")
			// Вызываем функцию чтения метрик из памяти и сохраняем их в хранилище.
			agent.ReadMemStats(&m)

		case <-reportTicker.C:

			logger.DebugLogger.Sugar().Debug("Вызываем отправку метрик")
			// Оборачиваем функцию ProcessBatch в функцию Retry, с попытками повтора в случае неудачи.
			fn := agent.Retry(agent.ProcessBatch, 3, 1*time.Second)
			// Пытаемся отправить данные на сервер.
			err := fn(context.Background(), cfg.ServerAddress, m)
			// Если отправка не удалась после всех попыток, логируем ошибку.
			if err != nil {
				logger.DebugLogger.Sugar().Error("Не удалось обработать пакет batch: ", err)
			}
		}
	}
}
