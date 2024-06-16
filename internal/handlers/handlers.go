package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	memstorage "github.com/RomanenkoDR/metrics/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
)

const (
	MyTypeGauge   string = "gauge"
	MyTypeCounter string = "counter"
)

func UpdateMetric(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	fmt.Printf("Получены параметры - тип: %s, имя: %s, значение: %s\n", metricType, metricName, metricValue)

	if metricName == "" {
		http.Error(res, "Имя метрики отсутствует", http.StatusBadRequest)
		return
	}

	// Проверка типа метрики
	if metricType != MyTypeGauge && metricType != MyTypeCounter {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Неверный тип метрики: %s", metricType)))
		return
	}

	switch metricType {
	case MyTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		storage.UpdateMetric(memstorage.MyTypeGauge, metricName, value)
		fmt.Printf("Получена метрика - тип: %s, имя: %s, значение: %f\n", metricType, metricName, value)
	case MyTypeCounter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		storage.UpdateMetric(memstorage.MyTypeCounter, metricName, value)
		fmt.Printf("Получена метрика - тип: %s, имя: %s, значение: %d\n", metricType, metricName, value)
	default:
		http.Error(res, fmt.Sprintf("Неверный тип метрики: %s", metricType), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Метрика обновлена"))
}

func ListMetrics(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)

	for name, value := range storage.GetAllMetrics() {
		fmt.Fprintf(res, "%s - %v\n", name, value)
	}
}

func GetValue(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	// Получение значения метрики из хранилища
	var value interface{}
	switch metricType {
	case MyTypeGauge:
		value = storage.GetGauge(metricName)
	case MyTypeCounter:
		value = storage.GetCounter(metricName)
	default:
		http.NotFound(res, req)
		return
	}

	// Проверка наличия метрики в хранилище
	if value == 0 {
		http.NotFound(res, req)
		return
	}

	// Возвращение значения метрики в текстовом виде
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("%v", value)))
}
