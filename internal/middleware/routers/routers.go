package routers

import (
	serverCfg "github.com/RomanenkoDR/metrics/internal/config/servercfg"
	handlersPcg "github.com/RomanenkoDR/metrics/internal/handlers"
	gzipPcg "github.com/RomanenkoDR/metrics/internal/middleware/gzip"
	loggerPcg "github.com/RomanenkoDR/metrics/internal/middleware/logger"

	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg serverCfg.Options) (chi.Router, handlersPcg.Handler, error) {
	h, err := handlersPcg.NewHandler(cfg.Filename, cfg.Restore)
	if err != nil {
		return nil, h, err
	}

	// Routers
	router := chi.NewRouter()
	router.Use(loggerPcg.LogHandler)
	router.Use(gzipPcg.GzipHandle)
	router.Get("/", h.HandleMain)
	router.Get("/ping", h.HandlePing)
	router.Post("/update/{type}/{metric}/{value}", h.HandleUpdate)
	router.Get("/value/gauge/{metric}", h.HandleValue)
	router.Get("/value/counter/{metric}", h.HandleValue)
	router.Post("/value/", h.HandleJSONValue)
	router.Post("/update/", h.HandleJSONUpdate)

	return router, h, nil
}
