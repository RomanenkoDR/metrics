package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/RomanenkoDR/metrics/internal/crypto"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// SetCryptoKey устанавливает путь к приватному ключу для расшифровки.
func (h *Handler) SetCryptoKey(path string) {
	h.PrivateKeyPath = path
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Обработка обновления метрики")
	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	switch metricType {
	case counterType:
		v, err := strconv.Atoi(value)
		if err != nil {
			logger.Error("Ошибка парсинга значения метрики", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Store.UpdateCounter(metric, storage.Counter(v))
	case gaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Error("Ошибка парсинга значения метрики", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Store.UpdateGauge(metric, storage.Gauge(v))
	default:
		logger.Warn("Некорректный тип метрики", zap.String("metricType", metricType))
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
}

// HandleUpdateJSON обрабатывает обновление одной метрики в формате JSON
func (h *Handler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var m Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("Ошибка чтения тела запроса", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := buf.Bytes()

	// Если передан ключ, расшифровываем AES-ключ, затем данные
	if h.PrivateKeyPath != "" {
		data, err = h.decryptPayload(data)
		if err != nil {
			http.Error(w, "Ошибка расшифровки данных", http.StatusBadRequest)
			return
		}
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		logger.Error("Ошибка десериализации JSON", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case counterType:
		if m.Delta == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.Store.UpdateCounter(m.ID, storage.Counter(*m.Delta))
	case gaugeType:
		if m.Value == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.Store.UpdateGauge(m.ID, storage.Gauge(*m.Value))
	default:
		logger.Warn("Некорректный тип метрики", zap.String("MType", m.MType))
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
	logger.Info("Метрика успешно обновлена", zap.Any("metric", m))
	w.WriteHeader(http.StatusOK)
}

// HandleUpdateBatch обрабатывает обновление нескольких метрик в формате JSON
func (h *Handler) HandleUpdateBatch(w http.ResponseWriter, r *http.Request) {
	var metrics []Metrics
	var buf bytes.Buffer

	// Читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Error("Ошибка чтения тела запроса", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := buf.Bytes()
	logger.Debug("Полученные данные для батча", zap.ByteString("raw_data", data))

	// Расшифровка, если передан ключ
	if h.PrivateKeyPath != "" {
		data, err = h.decryptPayload(data)
		if err != nil {
			http.Error(w, "Ошибка расшифровки данных", http.StatusBadRequest)
			return
		}
	}

	// Десериализуем JSON как массив метрик
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		logger.Error("Ошибка десериализации JSON (массив метрик)", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Логируем успешную десериализацию
	logger.Debug("Успешно десериализовано", zap.Int("batch_size", len(metrics)))

	// Обрабатываем каждую метрику
	for _, v := range metrics {
		switch v.MType {
		case counterType:
			if v.Delta == nil {
				http.Error(w, "metric value should not be empty", http.StatusBadRequest)
				return
			}
			h.Store.UpdateCounter(v.ID, storage.Counter(*v.Delta))
		case gaugeType:
			if v.Value == nil {
				http.Error(w, "metric value should not be empty", http.StatusBadRequest)
				return
			}
			h.Store.UpdateGauge(v.ID, storage.Gauge(*v.Value))
		default:
			logger.Warn("Некорректный тип метрики", zap.String("MType", v.MType))
			http.Error(w, "Incorrect metric type", http.StatusBadRequest)
			return
		}
	}

	logger.Info("Батч метрик успешно обработан", zap.Int("batch_size", len(metrics)))
	w.WriteHeader(http.StatusOK)
}

// decryptPayload расшифровывает полученные данные, если указан приватный ключ
func (h *Handler) decryptPayload(data []byte) ([]byte, error) {
	var encryptedPayload map[string][]byte
	err := json.Unmarshal(data, &encryptedPayload)
	if err != nil {
		logger.Error("Ошибка десериализации зашифрованного JSON", zap.Error(err))
		return nil, err
	}

	// Расшифровка AES-ключа
	aesKey, err := crypto.DecryptRSA(encryptedPayload["key"], h.PrivateKeyPath)
	if err != nil {
		logger.Error("Ошибка расшифровки AES-ключа", zap.Error(err))
		return nil, err
	}

	// Расшифровка данных
	decryptedData, err := crypto.DecryptAES(encryptedPayload["data"], aesKey)
	if err != nil {
		logger.Error("Ошибка расшифровки данных", zap.Error(err))
		return nil, err
	}

	return decryptedData, nil
}
