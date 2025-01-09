package server

import (
	"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
)

// setupStorage настраивает хранилище данных (файл или база данных).
func setupStorage(cfg types.OptionsServer) (storage.StorageWriter, error) {
	if cfg.DBDSN != "" {
		logger.DebugLogger.Info("Подключение к базе данных", zap.String("dsn", cfg.DBDSN))
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			logger.DebugLogger.Error("Ошибка подключения к базе данных", zap.Error(err))
			return nil, err
		}
		logger.DebugLogger.Info("Успешное подключение к базе данных")
		return database, nil // Возвращаем указатель, который реализует StorageWriter
	}

	logger.DebugLogger.Info("Используется локальное файловое хранилище", zap.String("file", cfg.Filename))
	return &storage.Localfile{Path: cfg.Filename}, nil
}
