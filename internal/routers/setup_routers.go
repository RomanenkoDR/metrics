package routers

import (
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func setupRoutes(router chi.Router, h handlers.Handler) {
	router.Get("/", h.HandleMain)
	router.Get("/ping", h.HandlePing)
	router.Get("/value/gauge/{metric}", h.HandleValue)
	router.Get("/value/counter/{metric}", h.HandleValue)

	router.Post("/update/{type}/{metric}/{value}", h.HandleUpdate)
	router.Post("/value/", h.HandleValueJSON)
	router.Post("/update/", h.HandleUpdateJSON)
	router.Post("/updates/", h.HandleUpdateBatch)
}
