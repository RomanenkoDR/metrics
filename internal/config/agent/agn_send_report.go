package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"net/http"
)

// Sender определяем тип функции, которая принимает контекст, строку с адресом сервера и объект MemStorage, и возвращает ошибку.
type sender func(context.Context, string, storage.MemStorage) error

// sendRequest - вспомогательная функция для отправки HTTP-запроса на сервер
func sendRequest(serverAddress string, data []byte) error {
	// Сжимаем данные перед отправкой на сервер
	compressedData, err := compress(data)
	if err != nil {
		return err
	}

	// Создание нового HTTP запроса типа POST с телом запроса в виде сжатого JSON
	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(compressedData))
	if err != nil {
		return err
	}

	// Устанавливаем заголовки запроса: тип контента, кодировка и поддержка сжатия
	request.Header.Set("Content-Type", contentTypeAppJSON)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

	// Создаем HTTP клиент для выполнения запроса
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	// Проверяем, успешно ли выполнен запрос (должен быть статус 200 OK)
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: %s; %s",
			"Can't send report to the server",
			resp.Status,
			b)
	}
	defer resp.Body.Close()
	return nil
}

// sendReport - функция для отправки одной метрики
func sendReport(serverAddress string, metrics Metrics) error {
	// Преобразование структуры метрики в JSON
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}

// sendReportBatch - функция для отправки нескольких метрик (батч)
func sendReportBatch(serverAddress string, metrics []Metrics) error {
	// Преобразование списка метрик в JSON
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	return sendRequest(serverAddress, data)
}
