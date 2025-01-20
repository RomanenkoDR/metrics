package agent

import (
	"github.com/RomanenkoDR/metrics/internal/storage"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// ReadMemStats обновляет метрики, используя пакет runtime
func ReadMemStats(m *storage.MemStorage, metricsCh chan storage.MemStorage) {
	var stat runtime.MemStats
	var mu sync.RWMutex

	mu.Lock()

	runtime.ReadMemStats(&stat)

	m.UpdateGauge("Alloc", storage.Gauge(stat.Alloc))
	m.UpdateGauge("BuckHashSys", storage.Gauge(stat.BuckHashSys))
	m.UpdateGauge("Frees", storage.Gauge(stat.Frees))
	m.UpdateGauge("GCCPUFraction", storage.Gauge(stat.GCCPUFraction))
	m.UpdateGauge("GCSys", storage.Gauge(stat.GCSys))
	m.UpdateGauge("HeapAlloc", storage.Gauge(stat.HeapAlloc))
	m.UpdateGauge("HeapIdle", storage.Gauge(stat.HeapIdle))
	m.UpdateGauge("HeapInuse", storage.Gauge(stat.HeapInuse))
	m.UpdateGauge("HeapObjects", storage.Gauge(stat.HeapObjects))
	m.UpdateGauge("HeapReleased", storage.Gauge(stat.HeapReleased))
	m.UpdateGauge("HeapSys", storage.Gauge(stat.HeapSys))
	m.UpdateGauge("LastGC", storage.Gauge(stat.LastGC))
	m.UpdateGauge("Lookups", storage.Gauge(stat.Lookups))
	m.UpdateGauge("MCacheInuse", storage.Gauge(stat.MCacheInuse))
	m.UpdateGauge("MCacheSys", storage.Gauge(stat.MCacheSys))
	m.UpdateGauge("MSpanInuse", storage.Gauge(stat.MSpanInuse))
	m.UpdateGauge("MSpanSys", storage.Gauge(stat.MSpanSys))
	m.UpdateGauge("Mallocs", storage.Gauge(stat.Mallocs))
	m.UpdateGauge("NextGC", storage.Gauge(stat.NextGC))
	m.UpdateGauge("NumForcedGC", storage.Gauge(stat.NumForcedGC))
	m.UpdateGauge("NumGC", storage.Gauge(stat.NumGC))
	m.UpdateGauge("OtherSys", storage.Gauge(stat.OtherSys))
	m.UpdateGauge("PauseTotalNs", storage.Gauge(stat.PauseTotalNs))
	m.UpdateGauge("StackInuse", storage.Gauge(stat.StackInuse))
	m.UpdateGauge("StackSys", storage.Gauge(stat.StackSys))
	m.UpdateGauge("Sys", storage.Gauge(stat.Sys))
	m.UpdateGauge("TotalAlloc", storage.Gauge(stat.TotalAlloc))
	m.UpdateGauge("RandomValue", storage.Gauge(rand.Float32()))
	m.UpdateCounter("PollCount", storage.Counter(1))

	vmem, _ := mem.VirtualMemory()
	cpu1, _ := cpu.Percent(time.Duration(0), true)

	m.UpdateGauge("TotalMemory", storage.Gauge(vmem.Total))
	m.UpdateGauge("FreeMemory", storage.Gauge(vmem.Free))
	m.UpdateGauge("CPUutilization1", storage.Gauge(cpu1[0]))

	mu.Unlock()

	metricsCh <- *m
}
