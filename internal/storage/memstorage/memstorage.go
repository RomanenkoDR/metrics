package memStorage

import "fmt"

type MetricType string

const (
	MyTypeGauge   MetricType = "gauge"
	MyTypeCounter MetricType = "counter"
)

// MemStorage - будет хранить метрики
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateMetric обновляет значение метрики в зависимости от её типа
func (m *MemStorage) UpdateMetric(metricType MetricType, name string, value interface{}) {
	switch metricType {
	case MyTypeGauge:
		if v, ok := value.(float64); ok {
			m.gauges[name] = v
			fmt.Println(metricType+" :", name+" =", v)
		}
	case MyTypeCounter:
		if v, ok := value.(int64); ok {
			m.counters[name] += v
			fmt.Println(metricType+" :", name+" =", m.counters[name])
		}
	}
}

// GetGauge возвращает значение метрики типа gauge
func (m *MemStorage) GetGauge(name string) float64 {
	value := m.gauges[name]
	return value
}

// GetCounter возвращает значение метрики типа counter
func (m *MemStorage) GetCounter(name string) int64 {
	value := m.counters[name]
	return value
}

func (m *MemStorage) GetAllMetrics() map[string]interface{} {
	allMetrics := make(map[string]interface{})
	for name, value := range m.gauges {
		allMetrics[name] = value
	}
	for name, value := range m.counters {
		allMetrics[name] = value
	}
	return allMetrics
}
