package mem

import "fmt"

type (
	Counter int64
	Gauge   float64
)

// MemStorage представляет собой хранилище метрик в памяти
type MemStorage struct {
	Data map[string]interface{}
}

// New создает новое хранилище метрик
func New() MemStorage {
	return MemStorage{
		Data: map[string]interface{}{},
	}
}

// Get возвращает значение метрики по имени
func (m *MemStorage) Get(metric string) (interface{}, error) {
	if v, ok := m.Data[metric]; ok {
		return v, nil
	}
	return "No such metric in memstorage", fmt.Errorf("metric not found")

}

// GetAll возвращает все метрики
func (m *MemStorage) GetAll() map[string]interface{} {
	return m.Data
}

// UpdateGauge обновляет значение метрики типа Gauge
func (m *MemStorage) UpdateGauge(metric string, value interface{}) {
	m.Data[metric] = value.(Gauge)
}

// UpdateCounter обновляет значение метрики типа Counter
func (m *MemStorage) UpdateCounter(metric string, value interface{}) {
	// Если метрика еще не существует, то создаем новую
	if m.Data[metric] == nil {
		m.Data[metric] = value.(Counter)
		return
	}
	// Если метрика существует, то увеличиваем ее значение
	m.Data[metric] = m.Data[metric].(Counter) + value.(Counter)
}
