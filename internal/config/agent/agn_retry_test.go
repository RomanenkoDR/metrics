package agent

import (
	"context"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/agent/types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	retryFunc := func(ctx context.Context, cfg types.OptionsAgent, metricsCh chan storage.MemStorage) error {
		return fmt.Errorf("Test error")
	}

	retriedFunc := Retry(retryFunc, 3, time.Millisecond*100)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := retriedFunc(ctx, types.OptionsAgent{}, make(chan storage.MemStorage))
	if err == nil {
		t.Error("Ожидалась ошибка, но её не было")
	}
}
