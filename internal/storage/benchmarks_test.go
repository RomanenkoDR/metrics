package storage

import (
	"context"
	"testing"
)

// MockWriter - реализация StorageWriter для тестирования
type MockWriter struct{}

func (mw *MockWriter) Write(s *MemStorage) error {
	return nil
}

func (mw *MockWriter) RestoreData(s *MemStorage) error {
	// Имитируем процесс восстановления данных
	s.CounterData["mock_counter"] = 100
	s.GaugeData["mock_gauge"] = 3.14
	return nil
}

func (mw *MockWriter) Save(ctx context.Context, t int, s *MemStorage) error {
	return nil
}

func (mw *MockWriter) Close() {}

// BenchmarkRestoreData тестирует производительность восстановления данных
func BenchmarkRestoreData(b *testing.B) {
	mockWriter := &MockWriter{}
	memStorage := New() // ✅ Теперь memStorage — это *MemStorage

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := mockWriter.RestoreData(memStorage) // ✅ Теперь передаётся *MemStorage
		if err != nil {
			b.Fatalf("RestoreData failed: %v", err)
		}
	}
}

func BenchmarkSaveData(b *testing.B) {
	mockWriter := &MockWriter{}
	memStorage := New() // ✅ Теперь memStorage — это *MemStorage
	memStorage.UpdateCounter("metric1", 42)
	memStorage.UpdateGauge("metric2", 3.14)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := mockWriter.Save(ctx, 10, memStorage) // ✅ Теперь передаётся *MemStorage
		if err != nil {
			b.Fatalf("SaveData failed: %v", err)
		}
	}
}
