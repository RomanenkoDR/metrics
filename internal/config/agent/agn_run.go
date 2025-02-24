package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
)

func Run() {
	// Логируем старт приложения
	logger.Info("Начало работы агента")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов: ", zap.Error(err))
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

	// Создаем контекст с отменой для управления завершением
	ctx, cancel := context.WithCancel(context.Background())

	// Канал для перехвата сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Запускаем горутину для обработки сигналов
	go func() {
		sig := <-sigChan
		logger.Info("Получен сигнал: ", zap.String("signal", sig.String()))
		cancel() // Отправляем сигнал завершения основному циклу
	}()

	// Запускаем основной цикл
loop:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Завершение работы агента. Отправка оставшихся данных...")
			flushData(cfg.ServerAddress, cfg.CryptoKey, &memStorage)
			logger.Info("Все данные успешно отправлены. Агент завершает работу.")
			break loop

		case <-pollTicker.C:
			logger.Debug("Сбор метрик")
			ReadMemStats(&memStorage)

		case <-reportTicker.C:
			logger.Debug("Отправка метрик")
			send := Retry(func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
				return ProcessBatch(ctx, serverAddress, cfg.CryptoKey, m)
			}, 3, 1*time.Second)

			err := send(ctx, cfg.ServerAddress, memStorage)
			if err != nil {
				logger.Error("Не удалось обработать пакет метрик: ", zap.Error(err))
			} else {
				logger.Info("Метрики отправлены на сервер")
			}
		}
	}
}

// flushData отправляет все накопленные метрики перед завершением работы агента.
func flushData(serverAddress, cryptoKey string, memStorage *storage.MemStorage) {
	send := Retry(func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		return ProcessBatch(ctx, serverAddress, cryptoKey, m)
	}, 3, 1*time.Second)

	err := send(context.Background(), serverAddress, *memStorage)
	if err != nil {
		logger.Error("Ошибка при финальной отправке данных: ", zap.Error(err))
	} else {
		logger.Info("Финальные данные успешно отправлены на сервер.")
	}
}
