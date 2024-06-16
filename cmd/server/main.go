package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RomanenkoDR/metrics/internal/config"
	handlers "github.com/RomanenkoDR/metrics/internal/handlers"
	memStorage "github.com/RomanenkoDR/metrics/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	metricType  = "/{metricType}"
	metricName  = "/{metricName}"
	metricValue = "/{metricValue}"
)

func main() {
	cfg := config.NewServerConfig()
	cfg.Init()

	// Создаем новый экземпляр Мемсторедж
	storage := memStorage.NewMemStorage()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/update"+metricType+metricName+metricValue,
		func(res http.ResponseWriter, req *http.Request) {
			handlers.UpdateMetric(res, req, storage)
		})
	r.Get("/value"+metricType+metricName,
		func(res http.ResponseWriter, req *http.Request) {
			metricName := chi.URLParam(req, "metricName")
			value := storage.GetGauge(metricName)
			if value == 0 {
				http.NotFound(res, req)
				return
			}
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintf("%v", value)))
		})
	r.Get("/",
		func(res http.ResponseWriter, req *http.Request) {
			handlers.ListMetrics(res, req, storage)
		})

	log.Printf("Запуск веб-сервера на %s\n", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
