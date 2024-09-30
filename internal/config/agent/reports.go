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
	"time"
)

// Инициализация клиента HTTP с таймаутом в 10 секунд
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Функция для отправки HTTP-запроса на сервер с данными
// ctx - контекст для управления таймингом и отменой запроса
// serverAddress - адрес сервера, куда отправляется запрос
// data - данные для отправки, сжатые и упакованные в HTTP-запрос
func sendRequest(ctx context.Context, serverAddress string, data []byte) error {
	compressedData, err := compress(data)
	if err != nil {
		log.Printf("Failed to compress data: %v", err)
		return err
	}

	// Создание нового HTTP-запроса с контекстом
	request, err := http.NewRequestWithContext(ctx, "POST", serverAddress, bytes.NewBuffer(compressedData))
	if err != nil {
		return err
	}

	// Установка заголовков запроса
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

	// Выполнение HTTP-запроса
	resp, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверка успешности выполнения запроса
	if resp.StatusCode != http.StatusOK {
		if resp.ContentLength > 0 {
			b, _ := io.ReadAll(resp.Body)
			log.Printf("Server response: %s", string(b))
		}
		return fmt.Errorf("can't send report to the server: %s", resp.Status)
	}

	return nil
}

// Функция для отправки одного отчета на сервер
// ctx - контекст
// serverAddress - адрес сервера
// metrics - данные метрики для отправки
func sendReport(ctx context.Context, serverAddress string, metrics Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(ctx, serverAddress, data)
}

// Функция для отправки батча (пакета) метрик на сервер
// ctx - контекст
// serverAddress - адрес сервера
// metrics - список метрик для отправки
func sendReportBatch(ctx context.Context, serverAddress string, metrics []Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(ctx, serverAddress, data)
}

// ProcessReport - функция для обработки и отправки каждой метрики отдельно
// ctx - контекст
// serverAddress - адрес сервера
// m - объект хранилища данных, содержащий метрики
func ProcessReport(ctx context.Context, serverAddress string, m storage.MemStorage) error {
	var metricsList []Metrics

	for k, v := range m.CounterData {
		metricsList = append(metricsList, Metrics{ID: k, MType: counterType, Delta: v})
	}
	for k, v := range m.GaugeData {
		metricsList = append(metricsList, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	for _, metrics := range metricsList {
		log.Println(metrics)
		if err := sendReport(ctx, serverAddress, metrics); err != nil {
			log.Printf("Failed to send metric ID: %s, error: %v", metrics.ID, err)
			return err
		}
	}
	return nil
}

// ProcessBatch - функция для обработки и отправки метрик батчем (пакетом)
// ctx - контекст
// serverAddress - адрес сервера
// m - объект хранилища данных, содержащий метрики
func ProcessBatch(ctx context.Context, serverAddress string, m storage.MemStorage) error {
	var metrics []Metrics

	for k, v := range m.CounterData {
		metrics = append(metrics, Metrics{ID: k, MType: counterType, Delta: v})
	}

	for k, v := range m.GaugeData {
		metrics = append(metrics, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	err := sendReportBatch(ctx, serverAddress, metrics)
	if err != nil {
		log.Printf("Failed to send metrics batch: %v", err)
		return err
	}
	return nil
}
