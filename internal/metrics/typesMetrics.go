package metrics

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
func NewMetrics() *SystemMetrics {
	return &SystemMetrics{
		metrics: make(map[string]SystemMetric),
	}
}
