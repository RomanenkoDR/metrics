package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Sender определяет тип функции, которая принимает контекст, адрес сервера и хранилище метрик, возвращая ошибку.
type Sender func(context.Context, string, storage.MemStorage) error

var Encrypt bool // Флаг для указания необходимости шифрования
var Key []byte   // Ключ для шифрования

// Retry создает обертку над функцией Sender, обеспечивающую повторные попытки выполнения функции в случае неудачи.
// retries задает количество попыток, а delay - задержку между попытками.
func Retry(sender Sender, retries int, delay time.Duration) Sender {
	return func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, serverAddress, m)
			// Если ошибок нет или попытки исчерпаны, возвращаем результат
			if err == nil || r >= retries {
				logger.DebugLogger.Sugar().Infof("Количество повторов: %d", r)
				return err
			}

			// Логируем попытку и ждем перед повторной попыткой
			logger.DebugLogger.Sugar().Warnf("Ошибка отправки метрик: %v, повтор через %v", err, delay)

			delay += time.Second * 2 // Увеличиваем задержку после каждой попытки

			select {
			case <-time.After(delay):
			case <-ctx.Done(): // Завершение контекста
				logger.DebugLogger.Sugar().Error("Контекст завершен, остановка повторных попыток")
				return ctx.Err()
			}
		}
	}
}
