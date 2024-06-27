package main

import (
	"net/http"

	config "github.com/RomanenkoDR/metrics/internal/config/serverConfig"
	handlers "github.com/RomanenkoDR/metrics/internal/handlers"
	logging "github.com/RomanenkoDR/metrics/internal/logging"
	memStorage "github.com/RomanenkoDR/metrics/internal/storage/mem"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	// Инициализация конфигурации сервера.
	configuration := config.NewServerConfiguration()
	configuration.InitServerConfiguration()

	// Инициализация хранилища MemStorage
	storage := memStorage.NewMemStorage()

	// Настройка логгера zap
	logger, err := logging.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	// Создание и настройка роутера chi
	r := chi.NewRouter()
	r.Use(logging.LoggingMiddleware(logger)) // Используем middleware логирование

	// Обновленные маршруты
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", func(res http.ResponseWriter, req *http.Request) {
			handlers.UpdateMetric(res, req, storage)
		})
		r.Post("/{metricType}/{metricName}/", func(res http.ResponseWriter, req *http.Request) {
			handlers.UpdateMetric(res, req, storage)
		})
	})
	r.Get("/value/{metricType}/{metricName}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetValueByName(res, req, storage)
	})

	r.Post("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.ListAllMetrics(res, req, storage)
	})
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		handlers.ListAllMetrics(res, req, storage)
	})

	// Логирование и обработка ошибок
	logger.Info("Запуск сервера на", zap.String("address", configuration.Address))
	err = http.ListenAndServe(configuration.Address, r)
	if err != nil {
		logger.Fatal("Ошибка при запуске сервера", zap.Error(err))
	}
}
