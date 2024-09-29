package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Counter int64
type Gauge float64

var Data *MemStorage

type MemStorage struct {
	CounterData map[string]Counter
	GaugeData   map[string]Gauge
}

func New() MemStorage {
	return MemStorage{
		CounterData: map[string]Counter{},
		GaugeData:   map[string]Gauge{},
	}
}

func (m *MemStorage) Get(metric string) (interface{}, error) {
	if v, ok := m.CounterData[metric]; ok {
		return v, nil
	}
	if v, ok := m.GaugeData[metric]; ok {
		return v, nil
	}
	return "No such metric in memstorage", fmt.Errorf("metric not found")

}

func (m *MemStorage) GetAllCounters() map[string]Counter {
	return m.CounterData
}

func (m *MemStorage) GetAllGauges() map[string]Gauge {
	return m.GaugeData
}

func (m *MemStorage) UpdateGauge(metric string, value Gauge) {
	m.GaugeData[metric] = value
}

func (m *MemStorage) UpdateCounter(metric string, value Counter) {
	m.CounterData[metric] = m.CounterData[metric] + value
}

func cleanFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	return nil
}

func WriteDataToFile(filename string, store MemStorage) error {
	//clear file
	err := cleanFile(filename)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
