package routers

import (
	"github.com/RomanenkoDR/metrics/internal/config/server"
	"github.com/RomanenkoDR/metrics/internal/middleware/gzip"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/middleware/token"
	"github.com/go-chi/chi/v5"
)

func setupMiddleware(router chi.Router, cfg server.Options) {
	router.Use(logger.LogHandler)
	router.Use(gzip.GzipHandle)
	if cfg.Key != "" {
		router.Use(token.CheckReqSign(cfg.Key))
	}
}
