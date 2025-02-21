package server

import (
	"flag"
	"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"os"
)

func parseOptions() (types.Options, error) {
	var cfg types.Options

	// Чтение флага "-a" для задания адреса сервера и порта
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Add address and port in format <address>:<port>")

	// Чтение флага "-i" для задания интервала сохранения метрик в файл
	flag.IntVar(&cfg.Interval, "i", 300, "Saving metrics to file interval")

	// Чтение флага "-f" для задания пути к файлу, где будут храниться метрики
	flag.StringVar(&cfg.Filename, "f", "./metrics.json", "File path")

	// Чтение флага "-r" для задания опции восстановления метрик из файла
	flag.BoolVar(&cfg.Restore, "r", true, "Restore metrics value from file")

	// Чтение флака "-k" для задания токена JWT
	flag.StringVar(&cfg.Key, "k", "", "Token auth by JWT")

	// Чтение флага "-d" для задания строки подключения к базе данных
	flag.StringVar(&cfg.DBDSN, "d", "", "Connection string in Postgres format")

	// Чтение параметра командной строки для пути к файлу с публичным ключом
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "Path to the public key for encryption")

	flag.StringVar(&cfg.Config, "c", "", "Path to config file")

	// Парсинг флагов командной строки
	flag.Parse()

	// Получение значений из переменных окружения и их применение
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	// Загружаем JSON, если указан файл конфигурации через `-c`
	var jsonCfg *serverFileConfig
	if cfg.Config != "" {
		jsonCfg, err = loadConfigFromFile(cfg.Config)
		if err != nil {
			logger.Warn("Ошибка загрузки JSON-конфигурации", zap.String("file", cfg.Config), zap.Error(err))
		} else {
			logger.Info("Загружена конфигурация из JSON-файла", zap.String("file", cfg.Config))
		}
	}

	// Переменная окружения CONFIG имеет приоритет над флагом `-c`
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		jsonCfg, err = loadConfigFromFile(envConfig)
		if err != nil {
			logger.Warn("Ошибка загрузки конфигурации из ENV CONFIG", zap.String("file", envConfig), zap.Error(err))
		} else {
			logger.Info("Загружена конфигурация из ENV CONFIG", zap.String("file", envConfig))
		}
	}

	// Применяем значения из JSON, если они не переопределены флагами или переменными окружения
	if jsonCfg != nil {
		if cfg.Address == "localhost:8080" {
			cfg.Address = jsonCfg.Address
		}
		if cfg.Interval == 300 {
			cfg.Interval = int(jsonCfg.StoreInterval.Seconds())
		}
		if cfg.Filename == "./metrics.json" {
			cfg.Filename = jsonCfg.StoreFile
		}
		if cfg.Restore {
			cfg.Restore = jsonCfg.Restore
		}
		if cfg.DBDSN == "" {
			cfg.DBDSN = jsonCfg.DatabaseDSN
		}
		if cfg.CryptoKey == "" {
			cfg.CryptoKey = jsonCfg.CryptoKey
		}
	}

	return cfg, nil
}
