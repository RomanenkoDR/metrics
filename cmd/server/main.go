package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RomanenkoDR/metrics/internal/config"
	handlers "github.com/RomanenkoDR/metrics/internal/handlers"
	memStorage "github.com/RomanenkoDR/metrics/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Инициализация конфигурации сервера.
	cfg := config.NewServerConfiguration()
	cfg.Init()

	// Изменение конфигурации на основе переменной окружения
	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.Address = addr
	}

	// Инициализация хранилища MemStorage
	storage := memStorage.NewMemStorage()

	// Создание и настройка роутера chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Обработчики для обновления и получения метрик
	// Получение метрик типом POST
	r.Post("/update/{metricType}/{metricName}/{metricValue}",
		func(res http.ResponseWriter, req *http.Request) {
			handlers.UpdateMetric(res, req, storage)
		})

	// Отправка метрик типом GET
	r.Get("/value/{metricType}/{metricName}",
		func(res http.ResponseWriter, req *http.Request) {
			handlers.GetValue(res, req, storage)
		})

	// Отправка ВСЕХ метрик типом GET
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.ListMetrics(res, req, storage)
	})

	// Логирование и обработка ошибок
	log.Printf("Запуск веб-сервера на %s\n", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
