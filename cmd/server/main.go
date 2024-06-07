package main

import (
	"log"
	"net/http"

	handlers "github.com/RomanenkoDR/metrics/internal/handlers"
)

// http://localhost:8080/update/counter/someMetric/527

// curl -v -X POST 'http://localhost:8080/update/counter/someMetric/527'

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(handlers.PostUpdate+handlers.MetricType+handlers.MetricName+handlers.MetricValue, handlers.Gauge)

	log.Println("Запуск веб-сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
