package routers

import (
	"github.com/RomanenkoDR/metrics/internal/config/server"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func InitRouter(cfg server.Options, h handlers.Handler) (chi.Router, error) {
	router := chi.NewRouter()

	setupMiddleware(router, cfg)
	setupRoutes(router, h)

	return router, nil
}
