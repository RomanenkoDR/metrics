package agent

import (
	"flag"
	"github.com/caarlos0/env"
)

// Options содержит параметры конфигурации агента.
type Options struct {
	ServerAddress  string `env:"ADDRESS"`         // Адрес сервера для отправки метрик.
	PollInterval   int    `env:"POLL_INTERVAL"`   // Интервал сбора метрик в секундах.
	ReportInterval int    `env:"REPORT_INTERVAL"` // Интервал отправки метрик в секундах.
	Key            string `env:"KEY"`             // Ключ для аутентификации JWT.
}

// ParseOptions парсит параметры конфигурации из переменных окружения и флагов.
func ParseOptions() (Options, error) {
	var opt Options

	// Чтение параметров из командной строки
	flag.IntVar(&opt.PollInterval, "p", 2, "Частота сбора метрик в секундах")
	flag.IntVar(&opt.ReportInterval, "r", 10, "Частота отправки метрик в секундах")
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "Адрес сервера для отправки метрик")
	flag.StringVar(&opt.Key, "k", "", "JWT токен для аутентификации")

	flag.Parse()

	// Чтение параметров из переменных окружения
	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}
