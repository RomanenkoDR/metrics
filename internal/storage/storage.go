package storage

import (
	"context"
	"fmt"
	"sync"
)

type Counter int64
type Gauge float64

type MemStorage struct {
	mu          sync.RWMutex
	CounterData map[string]Counter
	GaugeData   map[string]Gauge
}

// Интерфейс для различных хранилищ данных
type StorageWriter interface {
	Write(s MemStorage) error
	RestoreData(s *MemStorage) error
	Save(ctx context.Context, t int, s *MemStorage) error
	Close()
}

// Создание нового хранилища данных в памяти
func New() *MemStorage {
	return &MemStorage{
		CounterData: make(map[string]Counter),
		GaugeData:   make(map[string]Gauge),
	}
}

// Получение метрики
func (m *MemStorage) Get(metric string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.CounterData[metric]; ok {
		return v, nil
	}
	if v, ok := m.GaugeData[metric]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("метрика не найдена")
}

// Получение всех счетчиков
func (m *MemStorage) GetAllCounters() map[string]Counter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]Counter)
	for k, v := range m.CounterData {
		result[k] = v
	}
	return result
}

// Получение всех Gauge-метрик
func (m *MemStorage) GetAllGauge() map[string]Gauge {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]Gauge)
	for k, v := range m.GaugeData {
		result[k] = v
	}
	return result
}

// Обновление Gauge-метрики
func (m *MemStorage) UpdateGauge(metric string, value Gauge) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GaugeData[metric] = value
}

// Обновление Counter-метрики (инкремент)
func (m *MemStorage) UpdateCounter(metric string, value Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CounterData[metric] += value
}
