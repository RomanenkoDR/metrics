package metrics

import (
	"math/rand"
	"runtime"
)

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MetricType string

// Структура содержащая тип метрики и ее значение в виде интерфейса
type SystemMetric struct {
	Type  MetricType // Тип структур
	Value interface{}
}

// Структура содержащая переменную для счетчика кол-ва сборов метрики и мапу структур Metric
type SystemMetrics struct {
	metricCollectionCounter int64                   // Счетчик кол-ва сбора метрик
	metrics                 map[string]SystemMetric // мапа структур Metric
}

// Создание нового экземпляра структуры Metrics
func NewMetrics() *SystemMetrics { // Указатель на тип Metrics
	return &SystemMetrics{ // Создаем новый экземпляр, но возвращаем указатель на него
		metrics: make(map[string]SystemMetric), // Создаем и возвращаем новую мапу, строка-ключ, метрика-значение
	}
}

// Создание коллекции с метриками
func (m *SystemMetrics) CollectionOfMetrics() {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.metrics["Alloc"] = SystemMetric{Gauge, float64(memStats.Alloc)}
	m.metrics["BuckHashSys"] = SystemMetric{Gauge, float64(memStats.BuckHashSys)}
	m.metrics["Frees"] = SystemMetric{Gauge, float64(memStats.Frees)}
	m.metrics["GCCPUFraction"] = SystemMetric{Gauge, memStats.GCCPUFraction}
	m.metrics["GCSys"] = SystemMetric{Gauge, float64(memStats.GCSys)}
	m.metrics["HeapAlloc"] = SystemMetric{Gauge, float64(memStats.HeapAlloc)}
	m.metrics["HeapIdle"] = SystemMetric{Gauge, float64(memStats.HeapIdle)}
	m.metrics["HeapInuse"] = SystemMetric{Gauge, float64(memStats.HeapInuse)}
	m.metrics["HeapObjects"] = SystemMetric{Gauge, float64(memStats.HeapObjects)}
	m.metrics["HeapReleased"] = SystemMetric{Gauge, float64(memStats.HeapReleased)}
	m.metrics["HeapSys"] = SystemMetric{Gauge, float64(memStats.HeapSys)}
	m.metrics["LastGC"] = SystemMetric{Gauge, float64(memStats.LastGC)}
	m.metrics["Lookups"] = SystemMetric{Gauge, float64(memStats.Lookups)}
	m.metrics["MCacheInuse"] = SystemMetric{Gauge, float64(memStats.MCacheInuse)}
	m.metrics["MCacheSys"] = SystemMetric{Gauge, float64(memStats.MCacheSys)}
	m.metrics["MSpanInuse"] = SystemMetric{Gauge, float64(memStats.MSpanInuse)}
	m.metrics["MSpanSys"] = SystemMetric{Gauge, float64(memStats.MSpanSys)}
	m.metrics["Mallocs"] = SystemMetric{Gauge, float64(memStats.Mallocs)}
	m.metrics["NextGC"] = SystemMetric{Gauge, float64(memStats.NextGC)}
	m.metrics["NumForcedGC"] = SystemMetric{Gauge, float64(memStats.NumForcedGC)}
	m.metrics["NumGC"] = SystemMetric{Gauge, float64(memStats.NumGC)}
	m.metrics["OtherSys"] = SystemMetric{Gauge, float64(memStats.OtherSys)}
	m.metrics["PauseTotalNs"] = SystemMetric{Gauge, float64(memStats.PauseTotalNs)}
	m.metrics["StackInuse"] = SystemMetric{Gauge, float64(memStats.StackInuse)}
	m.metrics["StackSys"] = SystemMetric{Gauge, float64(memStats.StackSys)}
	m.metrics["Sys"] = SystemMetric{Gauge, float64(memStats.Sys)}
	m.metrics["TotalAlloc"] = SystemMetric{Gauge, float64(memStats.TotalAlloc)}

	m.metricCollectionCounter++
	m.metrics["PollCount"] = SystemMetric{Counter, m.metricCollectionCounter}
	m.metrics["RandomValue"] = SystemMetric{Gauge, rand.Float64()}
}

// Создаем копию метрик для отправки, что бы внешние изменения не повлияли на оригинальные данные
func (m *SystemMetrics) GetMetrics() map[string]SystemMetric {
	// Создаем новую мапу с тем же размером, что и "m.metrics"
	copyMetrics := make(map[string]SystemMetric, len(m.metrics))
	// Перебираем все элементы в "m.metrics"
	for key, value := range m.metrics {
		// Для каждой пары ключ-значение в `m.metrics` делаем запись в `copyMetrics`
		copyMetrics[key] = value
	}
	// Возвращаем созданную копию метрик
	return copyMetrics
}
