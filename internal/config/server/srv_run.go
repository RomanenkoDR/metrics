package server

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RunServer запускает сервер с настройкой маршрутов и graceful shutdown.
func RunServer() {
	logger.DebugLogger.Info("Запуск сервера...")

	// Парсим параметры конфигурации
	cfg, err := ParseOptions()
	if err != nil {
		logger.DebugLogger.Fatal("Ошибка парсинга конфигурации", zap.Error(err))
	}

	logger.DebugLogger.Info("Параметры конфигурации", zap.Any("config", cfg))

	// Настройка хранилища
	store, err := setupStorage(cfg)
	if err != nil {
		logger.DebugLogger.Fatal("Ошибка настройки хранилища", zap.Error(err))
	}

	// Инициализация обработчиков
	h := handlers.NewHandler()

	// Восстановление данных
	if cfg.Restore {
		logger.DebugLogger.Info("Восстановление данных из хранилища")
		if err := store.RestoreData(&h.Store); err != nil {
			logger.DebugLogger.Warn("Ошибка восстановления данных", zap.Error(err))
		}
	}

	// Инициализация маршрутизатора
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		logger.DebugLogger.Fatal("Ошибка инициализации маршрутизатора", zap.Error(err))
	}

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Инициализация MemStorage для метрик
	metrics := storage.New()

	// Настройка graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	setupShutdown(ctx, cancel, server, store, &metrics)

	// Запуск серверного цикла
	logger.DebugLogger.Info("Запуск HTTP сервера", zap.String("address", cfg.Address))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.DebugLogger.Fatal("Ошибка запуска HTTP сервера", zap.Error(err))
	}

	logger.DebugLogger.Info("Сервер завершил работу")
}
