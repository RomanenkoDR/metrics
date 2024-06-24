package logging

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware логирует информацию о запросах и ответах
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем запись для обертки ResponseWriter
		ww := &wrappedResponseWriter{w, http.StatusOK, 0}

		// Выполняем следующий обработчик
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Логируем информацию о запросе и ответе
		logrus.WithFields(logrus.Fields{
			"uri":           r.RequestURI,
			"method":        r.Method,
			"status":        ww.statusCode,
			"response_size": ww.responseSize,
			"duration_ms":   duration.Milliseconds(),
		}).Info("handled request")
	})
}

// wrappedResponseWriter обертка для http.ResponseWriter для захвата статуса и размера ответа
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (ww *wrappedResponseWriter) WriteHeader(statusCode int) {
	ww.statusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

func (ww *wrappedResponseWriter) Write(b []byte) (int, error) {
	size, err := ww.ResponseWriter.Write(b)
	ww.responseSize += size
	return size, err
}
