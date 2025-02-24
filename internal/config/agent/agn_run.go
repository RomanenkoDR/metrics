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
	logger.Info("Запуск агента...")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов", zap.Error(err))
	}

	// Если задан ключ шифрования, включаем его
	if cfg.Key != "" {
		Encrypt = true
		Key = []byte(cfg.Key)
		logger.Info("Шифрование включено")
	}

	// Создаём контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для перехвата системных сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Создаём тикеры
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	// Инициализируем хранилище метрик
	memStorage := storage.New()
	logger.Info("Инициализация хранилища успешна. Начало работы")

	// Основной цикл
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Остановка агента...")
				return

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
					logger.Error("Не удалось обработать пакет метрик", zap.Error(err))
				} else {
					logger.Info("Метрики успешно отправлены на сервер")
				}
			}
		}
	}()

	// Ожидаем сигнал завершения
	sig := <-sigChan
	logger.Info("Получен сигнал завершения", zap.String("signal", sig.String()))

	// Завершаем контекст, останавливаем агент
	cancel()

	// Ожидаем завершения всех операций перед выходом
	time.Sleep(1 * time.Second)
	logger.Info("Агент остановлен")
}
