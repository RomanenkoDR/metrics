package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"strings"
)

// ProcessReport Обрабатываем все метрики и отправляем их по одной на сервер
func ProcessReport(serverAddress string, m storage.MemStorage) error {
	var metrics Metrics

	// Формируем адрес для отправки метрик
	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	// Отправляем каждую метрику типа counter на сервер
	for k, v := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: v}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	// Отправляем каждую метрику типа gauge на сервер
	for k, v := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: v}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessBatch Функция для отправки батча (пакета) метрик
func ProcessBatch(ctx context.Context, serverAddress string, m storage.MemStorage) error {
	var metrics []Metrics

	// Формируем адрес для батч-отправки метрик
	serverAddress = strings.Join([]string{"http:/", serverAddress, "updates/"}, "/")

	// Добавляем все метрики типа counter в список для отправки
	for k, v := range m.CounterData {
		metrics = append(metrics, Metrics{ID: k, MType: counterType, Delta: v})
	}

	// Добавляем все метрики типа gauge в список для отправки
	for k, v := range m.GaugeData {
		metrics = append(metrics, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	// Отправляем батч метрик на сервер
	err := sendReportBatch(serverAddress, metrics)
	if err != nil {
		return err
	}
	return nil
}
