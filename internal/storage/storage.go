package storage

import (
	"errors"
)

// Counter и Gauge представляют типы метрик.
type Counter int64
type Gauge float64

// MemStorage хранит данные метрик в оперативной памяти.
type MemStorage struct {
	CounterData map[string]Counter
	GaugeData   map[string]Gauge
}

// StorageWriter определяет интерфейс для работы с хранилищем.
type StorageWriter interface {
	Write(s MemStorage) error        // Записывает данные.
	RestoreData(s *MemStorage) error // Восстанавливает данные.
	Save(t int, s MemStorage) error  // Сохраняет данные с интервалом.
	Close()                          // Закрывает хранилище.
}

// New создает новый экземпляр MemStorage.
func New() MemStorage {
	return MemStorage{
		CounterData: make(map[string]Counter),
		GaugeData:   make(map[string]Gauge),
	}
}

// Get возвращает значение метрики по имени.
func (m *MemStorage) Get(metric string) (interface{}, error) {
	if v, ok := m.CounterData[metric]; ok {
		return v, nil
	}
	if v, ok := m.GaugeData[metric]; ok {
		return v, nil
	}
	return nil, errors.New("metric not found")
}

// GetAllCounters возвращает все метрики типа Counter.
func (m *MemStorage) GetAllCounters() map[string]Counter {
	return m.CounterData
}

// GetAllGauge возвращает все метрики типа Gauge.
func (m *MemStorage) GetAllGauge() map[string]Gauge {
	return m.GaugeData
}

// UpdateGauge обновляет значение метрики типа Gauge.
func (m *MemStorage) UpdateGauge(metric string, value Gauge) {
	m.GaugeData[metric] = value
}

// UpdateCounter обновляет значение метрики типа Counter.
func (m *MemStorage) UpdateCounter(metric string, value Counter) {
	m.CounterData[metric] += value
}
