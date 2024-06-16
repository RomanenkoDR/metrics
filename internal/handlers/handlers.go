package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	memstorage "github.com/RomanenkoDR/metrics/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
)

// Константы для типов метрик
const (
	MyTypeGauge   string = "gauge"   // Константа для метрик типа gauge
	MyTypeCounter string = "counter" // Константа для метрик типа counter
)

// Обновляет метрику в хранилище на основе данных, полученных в запросе к серверу
func UpdateMetric(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	// Извлекаем параметры из URL
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	// Выводим информацию о полученных параметрах
	fmt.Printf("Получены параметры - тип: %s, имя: %s, значение: %s\n", metricType, metricName, metricValue)

	// Проверка, что имя метрики не пустое
	if metricName == "" {
		http.Error(res, "Имя метрики отсутствует", http.StatusBadRequest)
		return
	}

	// Проверка типа метрики на соотвествие одному из типов
	if metricType != MyTypeGauge && metricType != MyTypeCounter {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Неверный тип метрики: %s", metricType)))
		return
	}

	// Обновление метрики в зависимости от ее типа
	switch metricType {
	case MyTypeGauge: // Если тип метрики Gauge
		// Преобразование значения метрики в float64
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		// Обновление метрики типа gauge
		storage.UpdateMetric(memstorage.MyTypeGauge, metricName, value)
		fmt.Printf("Получена метрика - тип: %s, имя: %s, значение: %f\n", metricType, metricName, value)
	case MyTypeCounter: // Если тип метрики Counter
		// Преобразование значения метрики в int64
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "Некорректное значение метрики", http.StatusBadRequest)
			return
		}
		// Обновление метрики типа counter
		storage.UpdateMetric(memstorage.MyTypeCounter, metricName, value)
		fmt.Printf("Получена метрика - тип: %s, имя: %s, значение: %d\n", metricType, metricName, value)
	default: // Если тип метрики не соответсвует ни одному из условий
		http.Error(res, fmt.Sprintf("Неверный тип метрики: %s", metricType), http.StatusBadRequest)
		return
	}
	// Установка статуса ответа HTTP 200 (OK)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Метрика обновлена"))
}

// Вывод списка всех метрик из хранилища memStorage
func ListMetrics(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	// Установка заголовка Content-Type для ответа
	res.Header().Set("Content-Type", "text/plain")
	// Установка статуса ответа HTTP 200 (OK)
	res.WriteHeader(http.StatusOK)

	// Перебор всех метрик и их значений по типу имя метрики - значение
	for name, value := range storage.GetAllMetrics() {
		// Запись имени и значения каждой метрики в тело ответа
		fmt.Fprintf(res, "%s - %v\n", name, value)
	}
}

func GetValue(res http.ResponseWriter, req *http.Request, storage *memstorage.MemStorage) {
	// Извлечение параметров metricType и metricName из URL запроса
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
		// Если указан неверный тип метрики, возвращаем ошибку
		http.NotFound(res, req)
		return
	}

	// // Проверка наличия метрики в хранилище
	// if value == 0 {
	// 	// Если значение метрики равно нулю, считаем, что метрика не найдена и возвращаем ошибку 404 Not Found
	// 	http.NotFound(res, req)
	// 	return
	// }

	// Возвращение значения метрики в текстовом виде
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("%v", value)))
}
