package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Sender определяет тип функции для отправки метрик.
//
// Аргументы:
//   - context.Context: Контекст выполнения для управления временем жизни операции.
//   - string: Адрес сервера для отправки метрик.
//   - storage.MemStorage: Хранилище метрик для отправки.
//
// Возвращает:
//   - error: Ошибка, возникшая в процессе отправки.
type sender func(context.Context, string, storage.MemStorage) error

var (
	Encrypt bool   // Флаг шифрования.
	Key     []byte // Ключ для шифрования.
)

// Retry выполняет повторные попытки вызова функции sender при ошибке.
//
// Аргументы:
//   - sender: Функция для отправки метрик.
//   - retries: Количество попыток повторения.
//   - delay: Задержка между попытками.
//
// Возвращает:
//   - Функцию того же типа, которая обрабатывает повторные попытки.
func Retry(sender sender, retries int, delay time.Duration) sender {
	return func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, serverAddress, m)
			if err == nil || r >= retries {
				logger.DebugLogger.Sugar().Infof("Кол-во повторных попыток %d", r)
				return err
			}

			logger.DebugLogger.Sugar().Warnf("Отправка метрик завершила ошибкой: %v, повторная попытка через %v", err, delay)

			delay += time.Second * 2 // Увеличиваем задержку на 2 секунды после каждой попытки.

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				logger.DebugLogger.Sugar().Error("Контекст отменен, остановка повторных попыток")
				return ctx.Err()
			}
		}
	}
}
