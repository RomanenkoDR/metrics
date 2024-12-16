package agent

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"net/http"
)

var Encrypt bool
var Key []byte

// sender определяем тип функции, которая принимает контекст, строку с адресом сервера и объект MemStorage, и возвращает ошибку.
type sender func(context.Context, options, chan storage.MemStorage) error

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
	request.Header.Set("Content-Type", contentType)
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
			"Отправка на сервер не выполнена",
			resp.Status,
			b)
	}
	defer resp.Body.Close()
	return nil
}

// sendReport - функция для отправки одной метрики
func sendReport(serverAddress string, metrics Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	data, err = compress(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return err
	}

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

// sendReportBatch - функция для отправки нескольких метрик (батч)
func sendBatchReport(cfg options, metrics []Metrics) error {
	var sha256sum string

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	// Init request
	request, err := http.NewRequest("POST", cfg.ServerAddress, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	// Encrypt data and set Header
	if cfg.Encrypt {
		h := hmac.New(sha256.New, cfg.KeyByte)
		h.Write(data)
		sha256sum = hex.EncodeToString(h.Sum(nil))
		request.Header.Set("HashSHA256", sha256sum)
	}

	data, err = compress(data)
	if err != nil {
		return err
	}

	// Redefine request content
	request.Body = io.NopCloser(bytes.NewBuffer(data))

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return err
	}

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
