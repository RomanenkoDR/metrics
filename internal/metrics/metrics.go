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
type Metric struct {
	Type  MetricType // Тип структур
	Value interface{}
}

// Структура содержащая переменную для счетчика кол-ва сборов метрики и мапу структур Metric
type Metrics struct {
	metricCollectionCounter int64             // Счетчик кол-ва сбора метрик
	metrics                 map[string]Metric // мапа структур Metric
}

// Создание нового экземпляра структуры Metrics
func NewMetrics() *Metrics { // Указатель на тип Metrics
	return &Metrics{ // Создаем новый экземпляр, но возвращаем указатель на него
		metrics: make(map[string]Metric), // Создаем и возвращаем новую мапу, строка-ключ, метрика-значение
	}
}

// Создание коллекции с метриками
func (m *Metrics) CollectionOfMetrics() {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.metrics["Alloc"] = Metric{Gauge, float64(memStats.Alloc)}
	m.metrics["BuckHashSys"] = Metric{Gauge, float64(memStats.BuckHashSys)}
	m.metrics["Frees"] = Metric{Gauge, float64(memStats.Frees)}
	m.metrics["GCCPUFraction"] = Metric{Gauge, memStats.GCCPUFraction}
	m.metrics["GCSys"] = Metric{Gauge, float64(memStats.GCSys)}
	m.metrics["HeapAlloc"] = Metric{Gauge, float64(memStats.HeapAlloc)}
	m.metrics["HeapIdle"] = Metric{Gauge, float64(memStats.HeapIdle)}
	m.metrics["HeapInuse"] = Metric{Gauge, float64(memStats.HeapInuse)}
	m.metrics["HeapObjects"] = Metric{Gauge, float64(memStats.HeapObjects)}
	m.metrics["HeapReleased"] = Metric{Gauge, float64(memStats.HeapReleased)}
	m.metrics["HeapSys"] = Metric{Gauge, float64(memStats.HeapSys)}
	m.metrics["LastGC"] = Metric{Gauge, float64(memStats.LastGC)}
	m.metrics["Lookups"] = Metric{Gauge, float64(memStats.Lookups)}
	m.metrics["MCacheInuse"] = Metric{Gauge, float64(memStats.MCacheInuse)}
	m.metrics["MCacheSys"] = Metric{Gauge, float64(memStats.MCacheSys)}
	m.metrics["MSpanInuse"] = Metric{Gauge, float64(memStats.MSpanInuse)}
	m.metrics["MSpanSys"] = Metric{Gauge, float64(memStats.MSpanSys)}
	m.metrics["Mallocs"] = Metric{Gauge, float64(memStats.Mallocs)}
	m.metrics["NextGC"] = Metric{Gauge, float64(memStats.NextGC)}
	m.metrics["NumForcedGC"] = Metric{Gauge, float64(memStats.NumForcedGC)}
	m.metrics["NumGC"] = Metric{Gauge, float64(memStats.NumGC)}
	m.metrics["OtherSys"] = Metric{Gauge, float64(memStats.OtherSys)}
	m.metrics["PauseTotalNs"] = Metric{Gauge, float64(memStats.PauseTotalNs)}
	m.metrics["StackInuse"] = Metric{Gauge, float64(memStats.StackInuse)}
	m.metrics["StackSys"] = Metric{Gauge, float64(memStats.StackSys)}
	m.metrics["Sys"] = Metric{Gauge, float64(memStats.Sys)}
	m.metrics["TotalAlloc"] = Metric{Gauge, float64(memStats.TotalAlloc)}

	m.metricCollectionCounter++
	m.metrics["PollCount"] = Metric{Counter, m.metricCollectionCounter}
	m.metrics["RandomValue"] = Metric{Gauge, rand.Float64()}
}

// Создаем копию метрик для отправки, что бы внешние изменения не повлияли на оригинальные данные
func (m *Metrics) GetMetrics() map[string]Metric {
	// Создаем новую мапу с тем же размером, что и "m.metrics"
	copyMetrics := make(map[string]Metric, len(m.metrics))
	// Перебираем все элементы в "m.metrics"
	for key, value := range m.metrics {
		// Для каждой пары ключ-значение в `m.metrics` делаем запись в `copyMetrics`
		copyMetrics[key] = value
	}
	// Возвращаем созданную копию метрик
	return copyMetrics
}
