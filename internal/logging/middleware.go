package logging

import (
	"net/http"

	"go.uber.org/zap"
)

const defaultLogLevel = "INFO"

var log *zap.Logger

// NewLogger создает новый экземпляр логгера с заданным уровнем логирования
func NewLogger(level string) (*zap.Logger, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	log = zl

	return zl, err
}

// Logger возвращает текущий экземпляр логгера
func Logger() *zap.Logger {
	if log != nil {
		return log
	}

	var err error
	log, err = NewLogger(defaultLogLevel)
	if err != nil {
		panic(err)
	}

	return log
}

// LoggingMiddleware создает middleware для логирования запросов
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()))
			next.ServeHTTP(w, r)
		})
	}
}
