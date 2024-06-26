package main

import (
	"net/http"

	config "github.com/RomanenkoDR/metrics/internal/config/serverConfig"
	handlers "github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/logging"
	memStorage "github.com/RomanenkoDR/metrics/internal/storage/memStorage"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	// Инициализация конфигурации сервера.
	configuration := config.NewServerConfiguration()
	configuration.InitServerConfiguration()

	// Инициализация хранилища MemStorage
	storage := memStorage.NewMemStorage()

	// Настройка логгера logrus
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// })
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Создание и настройка роутера chi
	r := chi.NewRouter()
	r.Use(logging.LoggingMiddleware) // Используем middleware логирование

	// Получение метрик типом POST
	r.Post("/update/{metricType}/{metricName}/{metricValue}",
		func(res http.ResponseWriter, req *http.Request) {
			handlers.UpdateMetric(res, req, storage)
		})

	// Отправка метрик типом GET
	r.Get("/value/{metricType}/{metricName}",
		func(res http.ResponseWriter, req *http.Request) {
			handlers.GetValueByName(res, req, storage)
		})

	// Отправка ВСЕХ метрик типом GET
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.ListAllMetrics(res, req, storage)
	})

	// Логирование и обработка ошибок
	logrus.Infof("Запуск сервера на: %s", configuration.Address)
	err := http.ListenAndServe(configuration.Address, r)
	if err != nil {
		logrus.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
