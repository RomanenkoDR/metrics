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
	logger.Info("Запуск агента...")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов: ", zap.Error(err))
	}

	// Контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Канал для обработки сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Создаем таймеры для сбора и отправки метрик
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	// Инициализируем хранилище метрик
	memStorage := storage.New()
	logger.Info("Инициализация хранилища успешна. Начало работы")

	// Канал завершения работы
	done := make(chan struct{})

	// Основной процесс в горутине
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				logger.Info("Получен сигнал завершения. Отправляем последние метрики...")

				// Отправляем накопленные метрики перед завершением
				err := ProcessBatch(context.Background(), cfg.ServerAddress, cfg.CryptoKey, memStorage)
				if err != nil {
					logger.Error("Ошибка при отправке метрик перед завершением", zap.Error(err))
				}

				logger.Info("Агент завершил работу корректно.")
				return

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
					logger.Error("Ошибка отправки метрик", zap.Error(err))
				} else {
					logger.Info("Метрики успешно отправлены на сервер")
				}
			}
		}
	}()

	// Ожидаем сигнал завершения
	sig := <-sigChan
	logger.Info("Получен сигнал", zap.String("signal", sig.String()))
	cancel() // Отправляем сигнал завершения в контекст

	// Ожидаем завершения горутины
	<-done
	logger.Info("Агент завершил работу.")
}
