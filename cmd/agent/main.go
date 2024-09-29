package main

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/agent"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Sender определяем тип функции, которая принимает контекст, строку с адресом сервера и объект MemStorage, и возвращает ошибку.
type Sender func(context.Context, string, storage.MemStorage) error

// Retry функция принимает другую функцию Sender, количество попыток retries и задержку delay, возвращает функцию того же типа,
// которая выполняет sender с попытками повторов в случае неудачи.
func Retry(sender Sender, retries int, delay time.Duration) Sender {
	// Возвращаем новую функцию, которая пытается выполнить sender.
	return func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, serverAddress, m)
			// Если ошибок нет или количество попыток исчерпано, логируем результат и возвращаем ошибку (если она была).
			if err == nil || r >= retries {
				logger.DebugLogger.Sugar().Infof("Completed retry cycle with %d attempts", r)
				return err
			}

			// Логируем сообщение о неудачной попытке и увеличиваем задержку перед следующей попыткой.
			logger.DebugLogger.Sugar().Warnf("Function call failed: %v, retrying in %v", err, delay)

			delay += time.Second * 2 // Увеличиваем задержку на 2 секунды после каждой попытки.

			// Ожидаем либо окончания задержки, либо завершения контекста.
			select {
			case <-time.After(delay):
			case <-ctx.Done(): // Если контекст завершён (например, программа была остановлена), возвращаем ошибку контекста.
				logger.DebugLogger.Sugar().Error("Context cancelled, stopping retries")
				return ctx.Err()
			}
		}
	}
}

func main() {
	// Логируем старт приложения.
	logger.DebugLogger.Sugar().Info("Starting the application")

	// Парсим параметры командной строки с помощью функции из пакета agent.
	cfg, err := agent.ParseOptions()
	if err != nil {
		// Если произошла ошибка при парсинге, логируем фатальную ошибку и завершаем программу.
		logger.DebugLogger.Sugar().Fatal("Failed to parse options: ", err)
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
	logger.DebugLogger.Sugar().Info("Initialization storage completed. Start main function")

	// Основной цикл программы, который работает вечно (пока не завершится).
	for {
		select {

		case <-pollTicker.C:
			logger.DebugLogger.Sugar().Debug("Polling memory stats")
			// Вызываем функцию чтения метрик из памяти и сохраняем их в хранилище.
			agent.ReadMemStats(&m)

		case <-reportTicker.C:

			logger.DebugLogger.Sugar().Debug("Reporting memory stats")
			// Оборачиваем функцию ProcessBatch в функцию Retry, с попытками повтора в случае неудачи.
			fn := Retry(agent.ProcessBatch, 3, 1*time.Second)
			// Пытаемся отправить данные на сервер.
			err := fn(context.Background(), cfg.ServerAddress, m)
			// Если отправка не удалась после всех попыток, логируем ошибку.
			if err != nil {
				logger.DebugLogger.Sugar().Error("Failed to process batch: ", err)
			}
		}
	}
}
