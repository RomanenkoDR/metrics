package handlers

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/crypto"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type EncryptedPayload struct {
	AESKey string `json:"aes_key"`
	Data   string `json:"data"`
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.LogHTTPRequest(r)

	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	switch metricType {
	case counterType:
		v, err := strconv.Atoi(value)
		if err != nil {
			logger.Error("Ошибка конвертации counter", zap.String("value", value), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
			return
		}
		h.Store.UpdateCounter(metric, storage.Counter(v))
		logger.Info("Обновлен counter", zap.String("metric", metric), zap.Int("value", v))

	case gaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Error("Ошибка конвертации gauge", zap.String("value", value), zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
			return
		}
		h.Store.UpdateGauge(metric, storage.Gauge(v))
		logger.Info("Обновлен gauge", zap.String("metric", metric), zap.Float64("value", v))

	default:
		logger.Error("Некорректный тип метрики", zap.String("metricType", metricType))
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.LogHTTPResponse(http.StatusOK, time.Since(start), 0)
}

func (h *Handler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var m Metrics
	var buf bytes.Buffer
	start := time.Now()

	logger.LogHTTPRequest(r)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("Ошибка чтения тела запроса", zap.Error(err))
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	requestData := buf.Bytes()
	privKey, err := crypto.LoadPrivateKey(crypto.PrivateKeyPath)
	if err != nil {
		logger.Error("Ошибка загрузки приватного ключа", zap.Error(err))
		http.Error(w, "Ошибка загрузки приватного ключа", http.StatusInternalServerError)
		logger.LogHTTPResponse(http.StatusInternalServerError, time.Since(start), 0)
		return
	}

	decryptedData, err := decryptAndParsePayload(requestData, privKey)
	if err != nil {
		http.Error(w, "Ошибка расшифровки данных", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	if err := json.Unmarshal(decryptedData, &m); err != nil {
		logger.Error("Ошибка парсинга расшифрованного JSON", zap.Error(err), zap.String("decrypted_json", string(decryptedData)))
		http.Error(w, "Ошибка парсинга расшифрованного JSON", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	switch m.MType {
	case counterType:
		if m.Delta == nil {
			logger.Error("Отсутствует значение Delta")
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
			return
		}
		h.Store.UpdateCounter(m.ID, storage.Counter(*m.Delta))
		logger.Info("Обновлен counter", zap.String("metric", m.ID), zap.Int64("value", *m.Delta))
	case gaugeType:
		if m.Value == nil {
			logger.Error("Отсутствует значение Value")
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
			return
		}
		h.Store.UpdateGauge(m.ID, storage.Gauge(*m.Value))
		logger.Info("Обновлен gauge", zap.String("metric", m.ID), zap.Float64("value", *m.Value))
	default:
		logger.Error("Некорректный тип метрики", zap.String("MType", m.MType))
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.LogHTTPResponse(http.StatusOK, time.Since(start), len(decryptedData))
}

func (h *Handler) HandleUpdateBatch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.LogHTTPRequest(r)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("Ошибка чтения тела запроса", zap.Error(err))
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	requestData := buf.Bytes()
	privKey, err := crypto.LoadPrivateKey(crypto.PrivateKeyPath)
	if err != nil {
		logger.Error("Ошибка загрузки приватного ключа", zap.Error(err))
		http.Error(w, "Ошибка загрузки приватного ключа", http.StatusInternalServerError)
		logger.LogHTTPResponse(http.StatusInternalServerError, time.Since(start), 0)
		return
	}

	decryptedData, err := decryptAndParsePayload(requestData, privKey)
	if err != nil {
		http.Error(w, "Ошибка расшифровки данных", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	var metrics []Metrics
	if err := json.Unmarshal(decryptedData, &metrics); err != nil {
		logger.Error("Ошибка парсинга расшифрованного JSON", zap.Error(err), zap.String("decrypted_json", string(decryptedData)))
		http.Error(w, "Ошибка парсинга расшифрованного JSON", http.StatusBadRequest)
		logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
		return
	}

	for _, v := range metrics {
		switch v.MType {
		case counterType:
			if v.Delta == nil {
				logger.Error("Отсутствует значение Delta", zap.String("metric", v.ID))
				http.Error(w, "Значение counter не должно быть пустым", http.StatusBadRequest)
				logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
				return
			}
			h.Store.UpdateCounter(v.ID, storage.Counter(*v.Delta))
			logger.Info("Обновлен counter", zap.String("metric", v.ID), zap.Int64("value", *v.Delta))
		case gaugeType:
			if v.Value == nil {
				logger.Error("Отсутствует значение Value", zap.String("metric", v.ID))
				http.Error(w, "Значение gauge не должно быть пустым", http.StatusBadRequest)
				logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
				return
			}
			h.Store.UpdateGauge(v.ID, storage.Gauge(*v.Value))
			logger.Info("Обновлен gauge", zap.String("metric", v.ID), zap.Float64("value", *v.Value))
		default:
			logger.Error("Некорректный тип метрики", zap.String("metric", v.ID))
			http.Error(w, "Некорректный тип метрики", http.StatusBadRequest)
			logger.LogHTTPResponse(http.StatusBadRequest, time.Since(start), 0)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("Обновление метрик завершено. Ожидаем следующий запрос")
}

func decryptAndParsePayload(requestData []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	var encryptedPayload EncryptedPayload

	if err := json.Unmarshal(requestData, &encryptedPayload); err != nil {
		logger.Error("Ошибка парсинга зашифрованного JSON", zap.Error(err), zap.ByteString("request_data", requestData))
		return nil, fmt.Errorf("ошибка парсинга зашифрованного JSON: %w", err)
	}

	encryptedAESKey, err := base64.StdEncoding.DecodeString(encryptedPayload.AESKey)
	if err != nil {
		logger.Error("Ошибка декодирования AES-ключа", zap.Error(err))
		return nil, fmt.Errorf("ошибка декодирования AES-ключа: %w", err)
	}

	aesKey, err := crypto.DecryptAESKeyRSA(encryptedAESKey, privKey)
	if err != nil {
		logger.Error("Ошибка расшифровки AES-ключа", zap.Error(err))
		return nil, fmt.Errorf("ошибка расшифровки AES-ключа: %w", err)
	}

	encryptedData, err := base64.StdEncoding.DecodeString(encryptedPayload.Data)
	if err != nil {
		logger.Error("Ошибка декодирования данных", zap.Error(err))
		return nil, fmt.Errorf("ошибка декодирования данных: %w", err)
	}

	decryptedData, err := crypto.DecryptData(encryptedData, aesKey)
	if err != nil {
		logger.Error("Ошибка расшифровки данных", zap.Error(err))
		return nil, fmt.Errorf("ошибка расшифровки данных: %w", err)
	}

	// 🚀 Декодирование из base64 в JSON, если это строка в кавычках
	if len(decryptedData) > 0 && decryptedData[0] == '"' {
		var base64Str string
		if err := json.Unmarshal(decryptedData, &base64Str); err != nil {
			return nil, fmt.Errorf("ошибка преобразования JSON-строки: %w", err)
		}

		decodedJSON, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			return nil, fmt.Errorf("ошибка декодирования base64 JSON: %w", err)
		}

		decryptedData = decodedJSON
	}
	return decryptedData, nil
}
