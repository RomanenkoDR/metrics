package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// Структура, хранящая информацию о статусе и размере HTTP-ответа.
	responseData struct {
		status int
		size   int
	}
	// Обертка для http.ResponseWriter, которая сохраняет дополнительные данные о ответе (статус и размер).
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Метод Write для loggingResponseWriter увеличивает размер ответа
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// Метод WriteHeader для loggingResponseWriter записывает статус ответа
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LogHandler является middleware, который логирует информацию о запросах
func LogHandler(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// Создание логгера
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		sugar := logger.Sugar()

		// Инициализация данных ответа
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		// Передача запроса следующему обработчику
		next.ServeHTTP(&lw, r)

		// Вычисление времени выполнения запроса
		duration := time.Since(start)

		// Логирование информации о запросе и ответе
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
