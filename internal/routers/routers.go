package routers

import (
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/gzip"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/models"
	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg models.ConfigServer, h handlers.Handler) (chi.Router, error) {
	// Init rout for server
	router := chi.NewRouter()
	router.Use(logger.LogHandler)
	router.Use(gzip.GzipHandle)

	// Get rout
	router.Get("/", h.HandleMain)
	router.Get("/ping", h.HandlePing)
	router.Get("/value/gauge/{metric}", h.HandleValue)
	router.Get("/value/counter/{metric}", h.HandleValue)

	// Post rout
	router.Post("/update/{type}/{metric}/{value}", h.HandleUpdate)
	router.Post("/value/", h.HandleValueJSON)
	router.Post("/update/", h.HandleUpdateJSON)
	router.Post("/updates/", h.HandleUpdateBatch)

	return router, nil
}
