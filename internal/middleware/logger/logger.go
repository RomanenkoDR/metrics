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

func init() {
	// Создаем директорию, если она не существует
	dir := "./"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err) // Не удается создать директорию, приложение не может продолжить работу
	}

	// Создаем файл, если он не существует
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // Не удается создать файл логов, приложение не может продолжить работу
	}

	// Создаем файловый и консольный энкодеры
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Настроим уровни логирования
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
	)

	// Создаем новый Logger
	DebugLogger = zap.New(core)

	// Добавляем синхронизацию, чтобы избежать потери логов
	defer DebugLogger.Sync()
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		DebugLogger.Error("Failed to write response", zap.Error(err))
	}
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LogHandler(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sugar := DebugLogger.Sugar()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

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
