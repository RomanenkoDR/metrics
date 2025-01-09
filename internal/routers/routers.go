package routers

import (
	"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/gzip"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/middleware/token"
	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg types.OptionsServer, h handlers.Handler) (chi.Router, error) {
	// Init rout for server
	router := chi.NewRouter()

	// Use router
	router.Use(logger.LogHandler)
	router.Use(gzip.GzipHandle)
	if cfg.Key != "" {
		router.Use(token.CheckReqSign(cfg.Key))
	}

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
