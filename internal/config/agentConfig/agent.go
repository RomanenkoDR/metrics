package agentConfig

import (
	"flag"

	"github.com/caarlos0/env"
)

// Содержит параметры конфигурации для агента:
// адрес сервера, интервал опроса и интервал отправки метрики.
// Значения этих параметров могут быть заданы через переменные окружения.
type Options struct {
	ServerAddress  string `env:"ADDRESS"`         // Адрес сервера
	PollInterval   int    `env:"POLL_INTERVAL"`   // Интервал опроса в секундах для сбора метрик
	ReportInterval int    `env:"REPORT_INTERVAL"` // Интервал отправки отчета на сервер в секундах
}

// Парсит параметры конфигурации из командной строки и переменных окружения.
func ParseOptions() (Options, error) {
	var option Options

	// Установка значений по умолчанию и описаний для опций командной строки
	flag.IntVar(&option.PollInterval, "p", 2,
		"Частота в секундах для сбора метрик")
	flag.IntVar(&option.ReportInterval, "r", 10,
		"Частота в секундах для отправки отчета на сервер")
	flag.StringVar(&option.ServerAddress, "a", "localhost:8080",
		"Адрес сервера для отправки метрик")
	flag.Parse()

	// Парсит значения из переменных окружения и перезаписывает
	// соответствующие поля структуры option
	err := env.Parse(&option)
	if err != nil {
		return option, err
	}

	return option, nil
}
