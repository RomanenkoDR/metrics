package storage

import (
	"testing"
)

// MockWriter - реализация StorageWriter для тестирования
type MockWriter struct{}

func (mw *MockWriter) Write(s MemStorage) error {
	return nil
}

func (mw *MockWriter) RestoreData(s *MemStorage) error {
	// Имитируем процесс восстановления данных
	s.CounterData["mock_counter"] = 100
	s.GaugeData["mock_gauge"] = 3.14
	return nil
}

func (mw *MockWriter) Save(t int, s MemStorage) error {
	return nil
}

func (mw *MockWriter) Close() {}

func BenchmarkRestoreData(b *testing.B) {
	// Инициализация хранилища и mock Writer
	mockWriter := &MockWriter{}
	memStorage := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := mockWriter.RestoreData(&memStorage)
		if err != nil {
			b.Fatalf("RestoreData failed: %v", err)
		}
	}
}

func BenchmarkSaveData(b *testing.B) {
	// Инициализация хранилища и mock Writer
	mockWriter := &MockWriter{}
	memStorage := New()
	memStorage.UpdateCounter("metric1", 42)
	memStorage.UpdateGauge("metric2", 3.14)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := SaveData(memStorage, mockWriter)
		if err != nil {
			b.Fatalf("SaveData failed: %v", err)
		}
	}
}
