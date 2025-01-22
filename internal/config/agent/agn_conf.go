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
//
// Возвращает:
//   - Options: Структура с параметрами конфигурации.
//   - error: Ошибка, если не удалось прочитать параметры.
func ParseOptions() (Options, error) {
	var opt Options

	// Чтение параметров из командной строки

	flag.IntVar(&opt.PollInterval, "p", 2, "Частота сбора метрик в секундах")
	flag.IntVar(&opt.ReportInterval, "r", 10, "Частота отправки метрик в секундах")
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080", "Адрес сервера для отправки метрик")
	flag.StringVar(&opt.Key, "k", "", "JWT токен для аутентификации")

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
