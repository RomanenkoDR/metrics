package agent

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/crypto"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

// sendRequest - вспомогательная функция для отправки HTTP-запроса на сервер
func sendRequest(serverAddress string, data []byte, cryptoKeyPath string) error {
	// Генерируем случайный AES-ключ
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		logger.Error("Ошибка генерации AES-ключа", zap.Error(err))
		return err
	}

	// Шифруем данные с помощью AES
	encryptedData, err := crypto.EncryptAES(data, aesKey)
	if err != nil {
		logger.Error("Ошибка шифрования данных AES", zap.Error(err))
		return err
	}

	// Создаём JSON-объект для отправки
	payload := map[string][]byte{"data": encryptedData}

	// Если указан путь к RSA-ключу, шифруем и сам AES-ключ
	if cryptoKeyPath != "" {
		encryptedAESKey, err := crypto.EncryptRSA(aesKey, cryptoKeyPath)
		if err != nil {
			logger.Error("Ошибка шифрования AES-ключа RSA", zap.Error(err))
			return err
		}
		payload["key"] = encryptedAESKey
		logger.Info("Данные зашифрованы с использованием AES + RSA")
	} else {
		logger.Warn("Публичный ключ RSA не передан, используется только AES-шифрование")
	}

	// Сериализуем зашифрованные данные в JSON
	encryptedPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Ошибка сериализации зашифрованных данных", zap.Error(err))
		return err
	}

	// Создаём HTTP-запрос
	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(encryptedPayload))
	if err != nil {
		logger.Error("Ошибка создания HTTP-запроса", zap.Error(err))
		return err
	}

	// Устанавливаем заголовки
	request.Header.Set("Content-Type", "application/json")

	// Отправляем HTTP-запрос
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error("Ошибка выполнения HTTP-запроса", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа сервера
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		logger.Error("Ошибка при отправке метрик", zap.String("status", resp.Status), zap.String("response", string(b)))
		return fmt.Errorf("can't send report to the server: %s; %s", resp.Status, b)
	}

	logger.Info("Метрики успешно отправлены на сервер")
	return nil
}

// sendReport - отправка одной метрики
func sendReport(serverAddress, cryptoKeyPath string, metrics Metrics) error {
	logger.Debug("Подготовка к отправке метрики", zap.Any("metrics", metrics))
	data, err := json.Marshal(metrics)
	if err != nil {
		logger.Error("Ошибка сериализации метрики", zap.Error(err))
		return err
	}
	logger.Debug("Отправка метрики на сервер", zap.String("serverAddress", serverAddress))
	return sendRequest(serverAddress, data, cryptoKeyPath)
}

// sendReportBatch - отправка батча метрик
func sendReportBatch(serverAddress, cryptoKeyPath string, metrics []Metrics) error {
	logger.Debug("Подготовка к отправке батча метрик", zap.Int("batch_size", len(metrics)))
	data, err := json.Marshal(metrics)
	if err != nil {
		logger.Error("Ошибка сериализации батча метрик", zap.Error(err))
		return err
	}
	logger.Debug("Отправка батча метрик на сервер", zap.String("serverAddress", serverAddress))
	return sendRequest(serverAddress, data, cryptoKeyPath)
}

// ProcessReport - отправка метрик по одной
func ProcessReport(serverAddress, cryptoKeyPath string, m storage.MemStorage) error {
	var metrics Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	for k, v := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: v}
		logger.Debug("Отправка метрики", zap.Any("metrics", metrics))
		err := sendReport(serverAddress, cryptoKeyPath, metrics)
		if err != nil {
			logger.Error("Ошибка отправки метрики", zap.Error(err))
			return err
		}
	}

	for k, v := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: v}
		logger.Debug("Отправка метрики", zap.Any("metrics", metrics))
		err := sendReport(serverAddress, cryptoKeyPath, metrics)
		if err != nil {
			logger.Error("Ошибка отправки метрики", zap.Error(err))
			return err
		}
	}
	return nil
}

// ProcessBatch - отправка батча метрик
func ProcessBatch(ctx context.Context, serverAddress, cryptoKeyPath string, m storage.MemStorage) error {
	var metrics []Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "updates/"}, "/")

	for k, v := range m.CounterData {
		metrics = append(metrics, Metrics{ID: k, MType: counterType, Delta: v})
	}

	for k, v := range m.GaugeData {
		metrics = append(metrics, Metrics{ID: k, MType: gaugeType, Value: v})
	}

	return sendReportBatch(serverAddress, cryptoKeyPath, metrics)
}
