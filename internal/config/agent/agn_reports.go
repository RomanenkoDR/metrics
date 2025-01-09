package agent

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/agent/agn_types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// startReporting отправляет собранные метрики на сервер
func startReporting(ctx context.Context, cfg agn_types.OptionsAgent, metricsCh chan storage.MemStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка отправки метрик...")
			return
		case <-ticker.C:
			fn := Retry(ProcessBatch, 3, 1*time.Second)
			err := fn(ctx, cfg, metricsCh)
			if err != nil {
				log.Printf("Ошибка при отправке метрик: %v", err)
			}
		}
	}
}

// sendReport - вспомогательная функция для отправки HTTP-запроса на сервер
func sendReport(serverAddress string, metrics agn_types.Metrics) error {
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
	request.Header.Set("Content-Type", agn_types.ContentType)
	request.Header.Set("Content-Encoding", agn_types.Compression)
	request.Header.Set("Accept-Encoding", agn_types.Compression)

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

	var metrics agn_types.Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, vmem := range m.CounterData {
		metrics = agn_types.Metrics{ID: k, MType: agn_types.CounterType, Delta: vmem}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	for k, vmem := range m.GaugeData {
		metrics = agn_types.Metrics{ID: k, MType: agn_types.GaugeType, Value: vmem}
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendBatchReport(cfg agn_types.OptionsAgent, metrics []agn_types.Metrics) error {
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

	request.Header.Set("Content-Type", agn_types.ContentType)
	request.Header.Set("Content-Encoding", agn_types.Compression)
	request.Header.Set("Accept-Encoding", agn_types.Compression)

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

func ProcessBatch(ctx context.Context, cfg agn_types.OptionsAgent,
	metricsCh chan storage.MemStorage) error {
	var metrics []agn_types.Metrics

	// Receive MemStorage with actual metrics
	m := <-metricsCh

	// Prepare structure to send to the server
	for k, v := range m.CounterData {
		metrics = append(metrics, agn_types.Metrics{ID: k, MType: agn_types.CounterType, Delta: v})
	}
	for k, v := range m.GaugeData {
		metrics = append(metrics, agn_types.Metrics{ID: k, MType: agn_types.GaugeType, Value: v})
	}

	// Send report
	err := sendBatchReport(cfg, metrics)
	if err != nil {
		return err
	}

	return nil
}
