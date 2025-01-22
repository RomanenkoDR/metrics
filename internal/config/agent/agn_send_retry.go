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
	// Возвращаем новую функцию, которая пытается выполнить sender.
	return func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, serverAddress, m)
			// Если ошибок нет или количество попыток исчерпано, логируем результат и возвращаем ошибку (если она была).
			if err == nil || r >= retries {
				logger.DebugLogger.Sugar().Infof("Кол-во повторных попыток %d", r)
				return err
			}

			// Логируем сообщение о неудачной попытке и увеличиваем задержку перед следующей попыткой.
			logger.DebugLogger.Sugar().Warnf("Отправка метрик завершила ошибкой: %v, повторная попытка %v", err, delay)

			delay += time.Second * 2 // Увеличиваем задержку на 2 секунды после каждой попытки.

			// Ожидаем либо окончания задержки, либо завершения контекста.
			select {
			case <-time.After(delay):
			case <-ctx.Done(): // Если контекст завершён (например, программа была остановлена), возвращаем ошибку контекста.
				logger.DebugLogger.Sugar().Error("Context отменен, остановка повторных попыток")
				return ctx.Err()
			}
		}
	}
}
