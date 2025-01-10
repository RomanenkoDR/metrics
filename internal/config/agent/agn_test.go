package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"sync"
	"testing"
	"time"
)

func TestStartCollecting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	metrics := storage.New()
	metricsCh := make(chan storage.MemStorage, 1)

	go startCollecting(ctx, metrics, metricsCh, time.Millisecond*500)

	select {
	case <-ctx.Done():
		t.Error("Время ожидания истекло до получения метрик")
	case m := <-metricsCh:
		if len(m.GaugeData) == 0 && len(m.CounterData) == 0 {
			t.Error("Метрики не были собраны")
		}
	}
}

func TestStartSystemMetricsCollecting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	metrics := storage.New()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		startSystemMetricsCollecting(ctx, metrics, time.Millisecond*500)
	}()

	wg.Wait()

	if _, ok := metrics.GaugeData["TotalMemory"]; !ok {
		t.Error("Системные метрики не были собраны (TotalMemory отсутствует)")
	}

	if _, ok := metrics.GaugeData["FreeMemory"]; !ok {
		t.Error("Системные метрики не были собраны (FreeMemory отсутствует)")
	}
}
