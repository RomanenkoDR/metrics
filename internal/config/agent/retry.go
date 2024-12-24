package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"time"
)

type Sender func(context.Context, Options, chan storage.MemStorage) error

func Retry(sender Sender, retries int, delay time.Duration) Sender {
	return func(ctx context.Context, cfg Options, metricsCh chan storage.MemStorage) error {
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
