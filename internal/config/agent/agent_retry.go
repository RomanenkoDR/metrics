package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"time"
)

// Retry функция принимает другую функцию Sender, количество попыток retries и задержку delay,
// возвращает функцию того же типа,
// которая выполняет sender с попытками повторов в случае неудачи.
func Retry(sender sender, retries int, delay time.Duration) sender {
	return func(ctx context.Context, cfg options, metricsCh chan storage.MemStorage) error {
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
