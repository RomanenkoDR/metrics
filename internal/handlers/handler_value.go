package handlers

import (
	"encoding/json"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// HandleValue возвращает значение метрики по URL-параметру
func (h *Handler) HandleValue(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	v, err := h.Store.Get(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(v.(string))) // Приводим к строке
}

// HandleValueJSON возвращает значение метрики в формате JSON
func (h *Handler) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	var m Metrics

	// Логируем получение запроса
	logger.Debug("Получен запрос /value/", zap.String("URL", r.URL.Path))

	// Декодируем JSON-запрос
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		logger.Error("Ошибка декодирования JSON", zap.Error(err))
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Логируем распарсенный JSON
	logger.Debug("Распарсенные данные", zap.Any("metric", m))

	// Проверяем, существует ли метрика в хранилище
	var resp Metrics
	switch m.MType {
	case "counter":
		v, ok := h.Store.CounterData[m.ID]
		if !ok {
			logger.Warn("Метрика не найдена", zap.String("metric_id", m.ID))
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		vPtr := int64(v)
		resp = Metrics{
			ID:    m.ID,
			MType: "counter",
			Delta: &vPtr,
		}
	case "gauge":
		v, ok := h.Store.GaugeData[m.ID]
		if !ok {
			logger.Warn("Метрика не найдена", zap.String("metric_id", m.ID))
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		vPtr := float64(v)
		resp = Metrics{
			ID:    m.ID,
			MType: "gauge",
			Value: &vPtr,
		}
	default:
		logger.Error("Неизвестный тип метрики", zap.String("mType", m.MType))
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	// Кодируем ответ в JSON
	respJSON, err := json.Marshal(resp)
	if err != nil {
		logger.Error("Ошибка сериализации JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)

	// Логируем успешную обработку
	logger.Info("Метрика успешно получена", zap.String("metric_id", m.ID), zap.Any("value", resp))
}
