package handlers

import (
	"fmt"
	"net/http"

	memstorage "github.com/RomanenkoDR/metrics/internal/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

// Вывод значения метрики по параметрам запроса
func GetValue(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	// Извлечение параметров metricType и metricName из URL запроса
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	// Получение значения метрики из хранилища
	var value interface{}

	switch metricType {
	case Gauge:
		value = storage.GetGauge(metricName)
	case Counter:
		value = storage.GetCounter(metricName)
	default:
		http.NotFound(res, req)
		return
	}

	// Возвращение значения метрики в текстовом виде
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("%v", value)))
}

// Вывод списка всех метрик из хранилища memStorage
func ListMetrics(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)

	// Перебор всех метрик и их значений по типу имя метрики - значение
	for name, value := range storage.GetAllMetrics() {
		fmt.Fprintf(res, "%s - %v\n", name, value)
	}
}
