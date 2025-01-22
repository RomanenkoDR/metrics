package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Sender определяем тип функции, которая принимает контекст, строку с адресом сервера и объект MemStorage, и возвращает ошибку.
type sender func(context.Context, string, storage.MemStorage) error

var Encrypt bool
var Key []byte

// Retry функция принимает другую функцию Sender, количество попыток retries и задержку delay, возвращает функцию того же типа,
// которая выполняет sender с попытками повторов в случае неудачи.
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
