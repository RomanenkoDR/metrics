package routers

import (
	gzipPcg "github.com/RomanenkoDR/metrics/internal/compress"
	handlersPcg "github.com/RomanenkoDR/metrics/internal/handlers"
	loggerPcg "github.com/RomanenkoDR/metrics/internal/logging"
	"github.com/go-chi/chi/v5"
)

// InitRouter инициализирует маршрутизатор и задает маршруты
func InitRouter() chi.Router {
	// Создание нового обработчика
	h := handlersPcg.NewHandler()

	// Создание нового маршрутизатора
	router := chi.NewRouter()

	// Добавление middleware для логирования
	router.Use(loggerPcg.LogHandler)

	// Определение маршрутов
	router.Use(gzipPcg.CompressDataHandler)
	router.Get("/", h.GetListAllMetrics)                               // Маршрут для получения списка всех метрик.
	router.Post("/update/{type}/{metric}/{value}", h.PostUpdateMetric) // Маршрут для обновления метрики по типу, имени и значению
	router.Get("/value/gauge/{metric}", h.GetValueByName)              // Маршрут для получения значения метрики типа gauge по имени
	router.Get("/value/counter/{metric}", h.GetValueByName)            // Маршрут для получения значения метрики типа counter по имени
	router.Post("/value/", h.PostValueByJSON)                          // Маршрут для получения значения метрики по JSON-запросу
	router.Post("/update/", h.PostUpdateJSON)                          // Маршрут для обновления метрики по JSON-запросу
	return router
}
