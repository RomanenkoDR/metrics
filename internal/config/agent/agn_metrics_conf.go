package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"math/rand/v2"
	"runtime"
)

type Metrics struct {
	ID    string          `json:"id"`    // имя метрики
	MType string          `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta storage.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value storage.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}

const (
	contentTypeAppJSON string = "application/json"

	compression string = "gzip"

	counterType string = "counter"
	gaugeType   string = "gauge"
)

// ReadMemStats Renew metrics through runtime package
func ReadMemStats(m *storage.MemStorage) {

	var stat runtime.MemStats
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
}

// Функция для сжатия данных с использованием gzip
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	// Пишем исходные данные в gzip writer для сжатия
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	// Закрываем writer и завершаем процесс сжатия
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	// Возвращаем сжатые данные в виде байтового среза
	return b.Bytes(), nil
}
