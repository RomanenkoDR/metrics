package agent

import (
	"flag"
	"github.com/caarlos0/env"
)

// Options содержит конфигурационные параметры агента, которые могут быть заданы
// через переменные окружения или флаги командной строки.
type Options struct {
	ServerAddress  string `env:"ADDRESS"`         // Адрес сервера для отправки метрик
	PollInterval   int    `env:"POLL_INTERVAL"`   // Интервал сбора метрик (в секундах)
	ReportInterval int    `env:"REPORT_INTERVAL"` // Интервал отправки метрик (в секундах)
	Key            string `env:"KEY"`             // JWT токен для аутентификации
}

// ParseOptions парсит параметры конфигурации агента из флагов командной строки
// и переменных окружения. Возвращает структуру Options и ошибку, если произошла ошибка.
func ParseOptions() (Options, error) {
	var opt Options

	// Парсинг флагов командной строки
	flag.IntVar(&opt.PollInterval, "p", 2, "Frequency in seconds for collecting metrics")
	flag.IntVar(&opt.ReportInterval, "r", 10, "Frequency in seconds for sending report to the server")
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "Address of the server to send metrics")
	flag.StringVar(&opt.Key, "k", "", "Token auth by JWT")
	flag.Parse()

	// Парсинг переменных окружения
	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}
