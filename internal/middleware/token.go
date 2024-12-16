package middleware

import (
	"github.com/RomanenkoDR/metrics/internal/services"
	"net/http"
)

// AuthMiddleware проверяет JWT токен, используя функцию СheckAuthToken.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, statusCode, err := services.СheckAuthToken(r)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// Добавляем username как заголовок для последующих обработчиков
		r.Header.Set("X-Username", username)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
