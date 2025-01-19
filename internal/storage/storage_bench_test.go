package storage

import (
	"testing"
)

func BenchmarkMemStorage_UpdateGauge(b *testing.B) {
	storage := New() // New() возвращает *MemStorage

	b.ResetTimer() // Сбрасываем таймер для точного измерения
	for i := 0; i < b.N; i++ {
		storage.UpdateGauge("test_metric", Gauge(i*1))
	}
}

func BenchmarkMemStorage_UpdateCounter(b *testing.B) {
	storage := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.UpdateCounter("test_metric", Counter(i))
	}
}

func BenchmarkMemStorage_Get(b *testing.B) {
	storage := New()
	storage.UpdateGauge("test_metric", 100.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Get("test_metric")
		if err != nil {
			b.Fatalf("Ошибка получения метрики: %v", err)
		}
	}
}

func BenchmarkMemStorage_Write(b *testing.B) {
	storage := New()
	storage.UpdateGauge("test_metric", 100.5)
	storage.UpdateCounter("counter_metric", 42)

	localFile := &Localfile{Path: "test_metrics.json"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := localFile.Write(storage)
		if err != nil {
			b.Fatalf("Ошибка записи метрик: %v", err)
		}
	}
}

func BenchmarkMemStorage_RestoreData(b *testing.B) {
	storage := New()
	localFile := &Localfile{Path: "test_metrics.json"}

	// Подготавливаем данные для восстановления
	err := localFile.Write(storage)
	if err != nil {
		b.Fatalf("Ошибка подготовки файла для восстановления: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := localFile.RestoreData(&storage)
		if err != nil {
			b.Fatalf("Ошибка восстановления данных: %v", err)
		}
	}
}
