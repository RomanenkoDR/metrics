package agent

import (
	"flag"
	"github.com/caarlos0/env"
)

type Options struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
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

	// Флаг для пути к публичному ключу
	flag.StringVar(&opt.CryptoKey, "crypto-key", "", "Path to public key file")

	// Парсинг аргументов командной строки
	flag.Parse()

	// Парсинг переменных окружения и их присвоение в структуру Options
	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	// Возвращаем структуру с параметрами и nil (ошибки нет)
	return opt, nil
}
