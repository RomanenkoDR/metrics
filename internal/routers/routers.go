package routers

import (
	handlersPcg "github.com/RomanenkoDR/metrics/internal/handlers"
	loggerPcg "github.com/RomanenkoDR/metrics/internal/logging"
	"github.com/go-chi/chi/v5"
)

func InitRouter() chi.Router {
	h := handlersPcg.NewHandler()

	// Routers
	router := chi.NewRouter()
	router.Use(loggerPcg.LogHandler)
	router.Get("/", h.GetListAllMetrics)
	router.Post("/update/{type}/{metric}/{value}", h.PostUpdateMetric)
	router.Get("/value/gauge/{metric}", h.GetValueByName)
	router.Get("/value/counter/{metric}", h.GetValueByName)
	router.Post("/value/", h.PostValueByJSON)
	router.Post("/update/", h.PostUpdateJSON)
	return router
}
