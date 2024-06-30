package handlers

// Структура для метрик, используемая для сериализации/десериализации JSON-запросов и ответов.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

const (
	Counter = "counter"
	Gauge   = "gauge"
)
