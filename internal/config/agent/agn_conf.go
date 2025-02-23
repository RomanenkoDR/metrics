package agent

import (
	"flag"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"os"
)

type Options struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	Config         string `env:"CONFIG"`
}

func ParseOptions() (Options, error) {
	var opt Options

	// Чтение параметра командной строки для интервала сбора метрик (по умолчанию 2 секунды)
	flag.IntVar(&opt.PollInterval, "p", 2, "Frequency in seconds for collecting metrics")

	// Чтение параметра командной строки для интервала отправки метрик (по умолчанию 10 секунд)
	flag.IntVar(&opt.ReportInterval, "r", 10, "Frequency in seconds for sending report to the server")

	// Чтение параметра командной строки для адреса сервера (по умолчанию "localhost:8080")
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "Address of the server to send metrics")

	// Чтение параметра командной строки для установки JWT токена
	flag.StringVar(&opt.Key, "k", "", "Token auth by JWT")

	// Чтение параметра командной строки для пути к файлу с публичным ключом
	flag.StringVar(&opt.CryptoKey, "crypto-key", "", "Path to the public key for encryption")

	flag.StringVar(&opt.Config, "c", "", "Path to config file")

	// Парсинг аргументов командной строки
	flag.Parse()

	// Парсинг переменных окружения и их присвоение в структуру Options
	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	// Загружаем JSON, если указан файл конфигурации через `-c`
	var cfg *agentConfFromFile
	if opt.Config != "" {
		cfg, err = loadConfigFromFile(opt.Config)
		if err != nil {
			logger.Warn("Ошибка загрузки JSON-конфигурации", zap.Error(err))
		}
	}

	// Переменная окружения CONFIG имеет приоритет над флагом `-c`
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		cfg, err = loadConfigFromFile(envConfig)
		if err != nil {
			logger.Warn("Ошибка загрузки конфигурации из ENV CONFIG", zap.Error(err))
		}
	}

	// Применяем значения из JSON, если они не переопределены флагами или переменными окружения
	if cfg != nil {
		if opt.ServerAddress == "localhost:8080" {
			opt.ServerAddress = cfg.ServerAddress
		}
		if opt.PollInterval == 2 {
			opt.PollInterval = int(cfg.PollInterval.Seconds())
		}
		if opt.ReportInterval == 10 {
			opt.ReportInterval = int(cfg.ReportInterval.Seconds())
		}
		if opt.CryptoKey == "" {
			opt.CryptoKey = cfg.CryptoKey
		}
	}

	// Возвращаем структуру с параметрами и nil
	return opt, nil
}
