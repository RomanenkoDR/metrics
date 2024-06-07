package handlers

import (
	"net/http"
	"strconv"
	"strings"

	memstorage "github.com/RomanenkoDR/metrics/internal/memStorage"
)

const (
	PostUpdate  string = "/update"
	MetricType  string = "/{metricType}"
	MetricName  string = "/{metricName}"
	MetricValue string = "/{metricValue}"

	MyTypeGauge   string = "gauge"
	MyTypeCounter string = "counter"
)

var storage = memstorage.NewMemStorage()

func Gauge(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Метод отличен от POST"))
		//fmt.Println("Метод отличен от POST. Status code: ", http.StatusBadRequest)
		return
	}

	// Парсинг URL и получение значений метрик
	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != 5 {
		res.WriteHeader(http.StatusBadRequest)
		//res.Write([]byte("Некорректный URL запроса"))
		//fmt.Println("Некорректный URL. Status code: ", http.StatusBadRequest)
		return
	}
	metricType := urlParts[2]
	metricName := urlParts[3]
	metricValue := urlParts[4]

	if metricName == "" {
		res.WriteHeader(http.StatusBadRequest)
		//res.Write([]byte("Имя метрики отсутствует"))
		//fmt.Println("Имя метрики отсутствует. Status code: ", http.StatusBadRequest)
		return
	}

	// Проверка типа метрики
	if metricType != MyTypeGauge && metricType != MyTypeCounter {
		res.WriteHeader(http.StatusBadRequest)
		// res.Write([]byte(fmt.Sprintf("Неверный тип метрики: %s", metricType)))
		//fmt.Printf("Неверный тип метрики: %s. Status code: %d\n", metricType, http.StatusBadRequest)
		return
	}

	switch metricType {
	case MyTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateMetric(memstorage.MyTypeGauge, metricName, value)
	case MyTypeCounter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateMetric(memstorage.MyTypeCounter, metricName, value)
	}

	res.WriteHeader(http.StatusOK)
	//res.Write([]byte("POST запрос обработан"))
	//fmt.Println("POST запрос обработан. Status code: ", http.StatusOK)
}
