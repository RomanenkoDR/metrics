// agent.go

package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
)

// sendRequest отправляет HTTP-запрос на сервер.
//
// Аргументы:
//   - serverAddress: Адрес сервера для отправки данных.
//   - data: Данные в виде среза байт для отправки.
//
// Возвращает:
//   - error: Ошибка в процессе отправки, если произошла.
func sendRequest(serverAddress string, data []byte) error {
	compressedData, err := compress(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(compressedData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", contentTypeAppJSON)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("can't send report to the server: %s; %s", resp.Status, b)
	}

	return nil
}

// sendReport отправляет одну метрику на сервер.
//
// Аргументы:
//   - serverAddress: Адрес сервера для отправки данных.
//   - metrics: Метрика для отправки.
//
// Возвращает:
//   - error: Ошибка в процессе отправки, если произошла.
func sendReport(serverAddress string, metrics Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}

// sendReportBatch отправляет несколько метрик (батч) на сервер.
//
// Аргументы:
//   - serverAddress: Адрес сервера для отправки данных.
//   - metrics: Список метрик для отправки.
//
// Возвращает:
//   - error: Ошибка в процессе отправки, если произошла.
func sendReportBatch(serverAddress string, metrics []Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}

// ProcessReport обрабатывает все метрики и отправляет их по одной на сервер.
//
// Аргументы:
//   - serverAddress: Адрес сервера для отправки данных.
//   - m: Хранилище метрик для отправки.
//
// Возвращает:
//   - error: Ошибка в процессе обработки, если произошла.
func ProcessReport(serverAddress string, m storage.MemStorage) error {
	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	for k, v := range m.CounterData {
		metrics := Metrics{ID: k, MType: counterType, Delta: v}
		log.Println(metrics)
		if err := sendReport(serverAddress, metrics); err != nil {
			return err
		}
	}

	for k, v := range m.GaugeData {
		metrics := Metrics{ID: k, MType: gaugeType, Value: v}
		log.Println(metrics)
		if err := sendReport(serverAddress, metrics); err != nil {
			return err
		}
	}

	return nil
}

// ProcessBatch отправляет пакет метрик на сервер.
//
// Аргументы:
//   - ctx: Контекст выполнения.
//   - serverAddress: Адрес сервера для отправки данных.
//   - m: Хранилище метрик для отправки.
//
// Возвращает:
//   - error: Ошибка в процессе обработки, если произошла.
func ProcessBatch(ctx context.Context, serverAddress string, m storage.MemStorage) error {
	serverAddress = strings.Join([]string{"http:/", serverAddress, "updates/"}, "/")

	var metrics []Metrics

	for k, v := range m.CounterData {
		metrics = append(metrics, Metrics{ID: k, MType: counterType, Delta: v})
	}

	for k, v := range m.GaugeData {
		metrics = append(metrics, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	return sendReportBatch(serverAddress, metrics)
}
