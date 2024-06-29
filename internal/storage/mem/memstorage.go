package mem

import "fmt"

type Counter int64
type Gauge float64

type MemStorage struct {
	Data map[string]interface{}
}

func New() MemStorage {
	return MemStorage{
		Data: map[string]interface{}{},
	}
}

func (m *MemStorage) Get(metric string) (interface{}, error) {
	if v, ok := m.Data[metric]; ok {
		return v, nil
	}
	return "No such metric in memstorage", fmt.Errorf("metric not found")

}

func (m *MemStorage) GetAll() map[string]interface{} {
	return m.Data
}

func (m *MemStorage) UpdateGauge(metric string, value interface{}) {
	m.Data[metric] = value.(Gauge)
}

func (m *MemStorage) UpdateCounter(metric string, value interface{}) {
	if m.Data[metric] == nil {
		m.Data[metric] = value.(Counter)
		return
	}
	m.Data[metric] = m.Data[metric].(Counter) + value.(Counter)
}
