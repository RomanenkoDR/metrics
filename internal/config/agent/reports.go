package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"strings"
)

// ProcessReport Обрабатываем все метрики и отправляем их по одной на сервер
func ProcessReport(serverAddress string, m storage.MemStorage) error {
	// metric type variable

	var metrics Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, vmem := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: vmem}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	for k, vmem := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: vmem}
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessBatch Функция для отправки батча (пакета) метрик
func ProcessBatch(ctx context.Context, cfg options,
	metricsCh chan storage.MemStorage) error {
	var metrics []Metrics

	// Receive MemStorage with actual metrics
	m := <-metricsCh

	// Prepare structure to send to the server
	for k, v := range m.CounterData {
		metrics = append(metrics, Metrics{ID: k, MType: counterType, Delta: v})
	}
	for k, v := range m.GaugeData {
		metrics = append(metrics, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	// Send report
	err := sendBatchReport(cfg, metrics)
	if err != nil {
		return err
	}

	return nil
}
