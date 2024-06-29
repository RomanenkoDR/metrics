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
	router.Get("/", h.HandleMain)
	router.Post("/update/{type}/{metric}/{value}", h.HandleUpdate)
	router.Get("/value/gauge/{metric}", h.HandleValue)
	router.Get("/value/counter/{metric}", h.HandleValue)
	router.Post("/value/", h.HandleJSONValue)
	router.Post("/update/", h.HandleJSONUpdate)
	return router
}
