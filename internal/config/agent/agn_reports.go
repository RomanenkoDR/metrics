package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/crypto"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
)

// sendRequest - вспомогательная функция для отправки HTTP-запроса на сервер
func sendRequest(serverAddress string, data []byte) error {
	// Проверяем, есть ли публичный ключ
	if crypto.PublicKey != nil {
		encryptedData, err := crypto.EncryptData(data, crypto.PublicKey)
		if err != nil {
			log.Printf("Ошибка шифрования данных: %v", err)
			return fmt.Errorf("ошибка шифрования: %v", err)
		}
		data = encryptedData
	} else {
		log.Println("Публичный ключ не загружен, данные отправляются без шифрования")
	}

	// Сжимаем данные перед отправкой на сервер
	compressedData, err := compress(data)
	if err != nil {
		return fmt.Errorf("ошибка сжатия данных: %v", err)
	}

	// Создаём HTTP-запрос
	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(compressedData))
	if err != nil {
		return err
	}

	// Устанавливаем заголовки запроса
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка отправки отчёта: %s; %s", resp.Status, b)
	}

	return nil
}

// sendReport - функция для отправки одной метрики
func sendReport(serverAddress string, metrics Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}

// sendReportBatch - функция для отправки нескольких метрик
func sendReportBatch(serverAddress string, metrics []Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}

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
