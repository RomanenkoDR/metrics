package logger

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

const logfile string = "./info.log"

var (
	DebugLogger *zap.Logger
)

// Инициализация логгера
func init() {
	dir := "./"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
	)

	DebugLogger = zap.New(core)
	defer DebugLogger.Sync()
}

// Write записывает данные в HTTP-ответ и логирует ошибки записи, если они есть.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		DebugLogger.Error("Ошибка записи ответа", zap.Error(err))
	}
	r.responseData.size += size
	return size, err
}

// WriteHeader записывает HTTP-статус в ответ.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LogHandler добавляет логирование к обработке HTTP-запросов.
func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{status: 0, size: 0}
		lw := loggingResponseWriter{ResponseWriter: w, responseData: responseData}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		DebugLogger.Sugar().Infof("URI: %s, Метод: %s, Статус: %d, Длительность: %s, Размер: %d",
			r.RequestURI, r.Method, responseData.status, duration, responseData.size)
	})
}

// Debug записывает отладочное сообщение с дополнительным контекстом.
func Debug(message string, fields ...zap.Field) {
	DebugLogger.Debug(message, fields...)
}

// Info записывает информационное сообщение с дополнительным контекстом.
func Info(message string, fields ...zap.Field) {
	DebugLogger.Info(message, fields...)
}

// Warn записывает предупреждающее сообщение с дополнительным контекстом.
func Warn(message string, fields ...zap.Field) {
	DebugLogger.Warn(message, fields...)
}

// Error записывает сообщение об ошибке с дополнительным контекстом.
func Error(message string, fields ...zap.Field) {
	DebugLogger.Error(message, fields...)
}

// Fatal записывает сообщение о фатальной ошибке и завершает приложение.
func Fatal(message string, fields ...zap.Field) {
	DebugLogger.Fatal(message, fields...)
}

// LogHTTPRequest логирует детали HTTP-запроса для отладки.
func LogHTTPRequest(r *http.Request) {
	DebugLogger.Sugar().Infof("HTTP-запрос - Метод: %s, URI: %s, Заголовки: %v", r.Method, r.RequestURI, r.Header)
}

// LogHTTPResponse логирует детали HTTP-ответа для отладки.
func LogHTTPResponse(statusCode int, duration time.Duration, size int) {
	DebugLogger.Sugar().Infof("HTTP-ответ - Статус: %d, Длительность: %s, Размер: %d", statusCode, duration, size)
}
