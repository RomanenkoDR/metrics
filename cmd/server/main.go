package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	PostUpdate  string = "/update"
	MetricType  string = "/{metricType}"
	MetricName  string = "/{metricName}"
	MetricValue string = "/{metricValue}"

	MyTypeGauge   string = "gauge"
	MyTypeCounter string = "counter"
)

// http://localhost:8080/update/counter/someMetric/527

// curl -v -X POST 'http://localhost:8080/update/counter/someMetric/527'

func gauge(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Метод отличен от POST"))
		fmt.Println("Метод отличен от POST. Status code: ", http.StatusBadRequest)
		return
	}

	// Парсинг URL и получение значений метрик
	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != 5 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Некорректный URL запроса"))
		fmt.Println("Некорректный URL. Status code: ", http.StatusBadRequest)
		return
	}
	metricType := urlParts[2]
	metricName := urlParts[3]
	// metricValue := urlParts[4]

	if metricName == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Имя метрики отсутствует"))
		fmt.Println("Имя метрики отсутствует. Status code: ", http.StatusBadRequest)
		return
	}

	// Проверка типа метрики
	if metricType != MyTypeGauge && metricType != MyTypeCounter {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Неверный тип метрики: %s", metricType)))
		fmt.Printf("Неверный тип метрики: %s. Status code: %d\n", metricType, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("POST запрос обработан"))
	fmt.Println("POST запрос обработан. Status code: ", http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(PostUpdate+MetricType+MetricName+MetricValue, gauge)

	log.Println("Запуск веб-сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
