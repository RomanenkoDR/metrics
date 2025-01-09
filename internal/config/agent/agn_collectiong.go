package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"time"
)

// startCollecting собирает метрики из runtime и передает их в канал
func startCollecting(ctx context.Context, m storage.MemStorage, metricsCh chan storage.MemStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка сбора метрик runtime...")
			return
		case <-ticker.C:
			ReadMemStats(&m, metricsCh)
		}
	}
}

// startSystemMetricsCollecting собирает системные метрики через gopsutil
func startSystemMetricsCollecting(ctx context.Context, m storage.MemStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка сбора системных метрик...")
			return
		case <-ticker.C:
			vmem, _ := mem.VirtualMemory()
			cpuUtilization, _ := cpu.Percent(0, false)

			m.UpdateGauge("TotalMemory", storage.Gauge(vmem.Total))
			m.UpdateGauge("FreeMemory", storage.Gauge(vmem.Free))
			if len(cpuUtilization) > 0 {
				m.UpdateGauge("CPUutilization1", storage.Gauge(cpuUtilization[0]))
			}
		}
	}
}
