package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/agent/types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"time"
)

// Sender определяем тип функции, которая принимает контекст, строку с адресом сервера и объект MemStorage, и возвращает ошибку.
type Sender func(context.Context, types.OptionsAgent, chan storage.MemStorage) error

var Encrypt bool
var Key []byte

// Retry функция принимает другую функцию Sender, количество попыток retries и задержку delay, возвращает функцию того же типа,
// которая выполняет sender с попытками повторов в случае неудачи.
func Retry(sender Sender, retries int, delay time.Duration) Sender {
	return func(ctx context.Context, cfg types.OptionsAgent, metricsCh chan storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, cfg, metricsCh)
			if err == nil || r >= retries {
				// Return when there is no error or the maximum amount
				// of retries is reached.
				return err
			}

			log.Printf("Function call failed, retrying in %v", delay)

			// Increase delay
			delay = delay + time.Second*2

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
