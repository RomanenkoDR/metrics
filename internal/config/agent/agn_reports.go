package agent

import (
	"bytes"
	"context"
	"encoding/base64"
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

// sendRequest - отправляет данные на сервер (с шифрованием AES и RSA)
func sendRequest(serverURL string, data interface{}) error {
	logger.Info("Подготовка к отправке данных", zap.String("server_url", serverURL))

	// Кодируем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Ошибка сериализации данных в JSON", zap.Error(err))
		return err
	}

	// Генерируем AES-ключ
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		logger.Error("Ошибка генерации AES-ключа", zap.Error(err))
		return err
	}

	// Шифруем данные перед отправкой
	encryptedData, err := crypto.EncryptData(jsonData, aesKey)
	if err != nil {
		logger.Error("Ошибка шифрования данных", zap.Error(err))
		return err
	}

	// Загружаем публичный RSA-ключ
	pubKey, err := crypto.LoadPublicKey(crypto.PublicKeyPath)
	if err != nil {
		logger.Error("Ошибка загрузки публичного ключа", zap.Error(err))
		return err
	}

	// Шифруем AES-ключ с помощью публичного RSA-ключа
	encryptedAESKey, err := crypto.EncryptAESKeyRSA(aesKey, pubKey)
	if err != nil {
		logger.Error("Ошибка шифрования AES-ключа", zap.Error(err))
		return err
	}

	// Кодируем зашифрованные данные в base64
	encryptedDataB64 := base64.StdEncoding.EncodeToString(encryptedData)
	encryptedAESKeyB64 := base64.StdEncoding.EncodeToString(encryptedAESKey)

	// Формируем JSON
	payload := []byte(`{"aes_key":"` + encryptedAESKeyB64 + `", "data":"` + encryptedDataB64 + `"}`)

	// Создаём HTTP-запрос
	request, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Ошибка создания HTTP-запроса", zap.Error(err))
		return err
	}

	// Устанавливаем заголовки запроса
	request.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error("Ошибка выполнения HTTP-запроса", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Читаем тело ответа сервера
	respBody, _ := io.ReadAll(resp.Body)

	// Логируем статус ответа
	logger.Info("Получен ответ от сервера", zap.Int("status", resp.StatusCode), zap.String("response_body", string(respBody)))

	// Проверяем успешность ответа
	if resp.StatusCode != http.StatusOK {
		logger.Error("Ошибка при отправке данных", zap.Int("status", resp.StatusCode), zap.String("response", string(respBody)))
		return fmt.Errorf("ошибка при отправке данных: %d %s; ответ: %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(respBody))
	}

	logger.Info("Данные успешно зашифрованы и отправлены на сервер")
	return nil
}

func sendReport(serverAddress string, metrics Metrics) error {
	logger.Info("Отправка отчёта", zap.String("server", serverAddress), zap.Any("metrics", metrics))

	data, err := json.Marshal(metrics)
	if err != nil {
		logger.Error("Ошибка сериализации метрики", zap.Error(err))
		return err
	}

	err = sendRequest(serverAddress, data)
	if err != nil {
		logger.Error("Ошибка отправки отчёта", zap.String("server", serverAddress), zap.Error(err))
	}
	return err
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
	var hasErrors bool

	// Формируем адрес для отправки метрик
	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	// Отправляем каждую метрику типа counter
	for k, v := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: v}
		logger.Info("Отправка counter метрики", zap.String("id", k), zap.Int64("value", int64(v)))

		err := sendReport(serverAddress, metrics)
		if err != nil {
			logger.Error("Ошибка отправки counter метрики", zap.String("id", k), zap.Error(err))
			hasErrors = true
		}
	}

	// Отправляем каждую метрику типа gauge
	for k, v := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: v}
		logger.Info("Отправка gauge метрики", zap.String("id", k), zap.Float64("value", float64(v)))

		err := sendReport(serverAddress, metrics)
		if err != nil {
			logger.Error("Ошибка отправки gauge метрики", zap.String("id", k), zap.Error(err))
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("ошибка отправки некоторых метрик")
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
