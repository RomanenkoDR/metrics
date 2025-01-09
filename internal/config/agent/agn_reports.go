package agent

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/agent/types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
)

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
	request.Header.Set("Content-Type", types.ContentType)
	request.Header.Set("Content-Encoding", types.Compression)
	request.Header.Set("Accept-Encoding", types.Compression)

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
func sendReport(serverAddress string, metrics types.Metrics) error {
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
	request.Header.Set("Content-Type", types.ContentType)
	request.Header.Set("Content-Encoding", types.Compression)
	request.Header.Set("Accept-Encoding", types.Compression)

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

// Process all the metrics and send them to the server one by one
func ProcessReport(serverAddress string, m storage.MemStorage) error {
	// metric type variable

	var metrics types.Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, vmem := range m.CounterData {
		metrics = types.Metrics{ID: k, MType: types.CounterType, Delta: vmem}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	for k, vmem := range m.GaugeData {
		metrics = types.Metrics{ID: k, MType: types.GaugeType, Value: vmem}
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendBatchReport(cfg types.OptionsAgent, metrics []types.Metrics) error {
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

	request.Header.Set("Content-Type", types.ContentType)
	request.Header.Set("Content-Encoding", types.Compression)
	request.Header.Set("Accept-Encoding", types.Compression)

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

func ProcessBatch(ctx context.Context, cfg types.OptionsAgent,
	metricsCh chan storage.MemStorage) error {
	var metrics []types.Metrics

	// Receive MemStorage with actual metrics
	m := <-metricsCh

	// Prepare structure to send to the server
	for k, v := range m.CounterData {
		metrics = append(metrics, types.Metrics{ID: k, MType: types.CounterType, Delta: v})
	}
	for k, v := range m.GaugeData {
		metrics = append(metrics, types.Metrics{ID: k, MType: types.GaugeType, Value: v})
	}

	// Send report
	err := sendBatchReport(cfg, metrics)
	if err != nil {
		return err
	}

	return nil
}
