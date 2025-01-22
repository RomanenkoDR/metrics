package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"math/rand/v2"
	"runtime"
)

// Metrics представляет структуру метрики, которая используется для отправки на сервер.
type Metrics struct {
	ID    string          `json:"id"`    // Имя метрики
	MType string          `json:"type"`  // Тип метрики: gauge или counter
	Delta storage.Counter `json:"delta"` // Значение метрики для counter
	Value storage.Gauge   `json:"value"` // Значение метрики для gauge
}

const (
	contentTypeAppJSON string = "application/json" // Тип контента для JSON
	compression        string = "gzip"             // Метод сжатия данных
	counterType        string = "counter"          // Тип метрики counter
	gaugeType          string = "gauge"            // Тип метрики gauge
)

// ReadMemStats обновляет метрики агента, используя данные о состоянии памяти из пакета runtime.
func ReadMemStats(m *storage.MemStorage) {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)

	// Обновление gauge-метрик на основе статистики памяти
	m.UpdateGauge("Alloc", storage.Gauge(stat.Alloc))
	// ... (остальные метрики пропущены для краткости)
	m.UpdateGauge("RandomValue", storage.Gauge(rand.Float32()))
	m.UpdateCounter("PollCount", storage.Counter(1))
}

// compress сжимает данные с использованием алгоритма gzip и возвращает сжатые данные или ошибку.
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gzip writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to gzip writer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	return b.Bytes(), nil
}
