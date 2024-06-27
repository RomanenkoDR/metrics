package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	memStorage "github.com/RomanenkoDR/metrics/internal/storage/mem"
	"github.com/go-chi/chi/v5"
)

// Обновляет метрику в хранилище на основе данных, полученных в запросе к серверу
func UpdateMetric(res http.ResponseWriter, req *http.Request, storage *memStorage.MemStorage) {
	// Извлекаем параметры из URL
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	// Проверка, что имя метрики не пустое
	if metricName == "" {
		http.Error(res, "Имя метрики отсутствует", http.StatusBadRequest)
		return
	}

	// Проверка типа метрики на соотвествие одному из типов
	if metricType != Gauge && metricType != Counter {
		http.Error(res, "Некорректный тип метрики", http.StatusBadRequest)
		return
	}

	// Обновление метрики в зависимости от ее типа
	switch metricType {
	case Gauge: // Если тип метрики Gauge
		// Преобразование значения метрики в float64
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		// Обновление метрики типа gauge
		storage.UpdateMetric(memStorage.MyTypeGauge, metricName, value)
	case Counter: // Если тип метрики Counter
		// Преобразование значения метрики в int64
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		// Обновление метрики типа counter
		storage.UpdateMetric(memStorage.MyTypeCounter, metricName, value)
	default: // Если тип метрики не соответсвует ни одному из условий
		http.Error(res, fmt.Sprintf("Неверный тип метрики: %s", metricType), http.StatusBadRequest)
		return
	}
	// Установка статуса ответа HTTP 200 (OK)
	res.WriteHeader(http.StatusOK)
}
