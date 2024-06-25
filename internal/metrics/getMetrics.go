package metrics

// Создаем копию метрик для отправки, что бы внешние изменения не повлияли на оригинальные данные
func (m *SystemMetrics) GetMetrics() map[string]SystemMetric {
	// Создаем новую мапу с тем же размером, что и "m.metrics"
	copyMetrics := make(map[string]SystemMetric, len(m.metrics))
	for key, value := range m.metrics {
		copyMetrics[key] = value
	}
	// Возвращаем созданную копию метрик
	return copyMetrics
}
