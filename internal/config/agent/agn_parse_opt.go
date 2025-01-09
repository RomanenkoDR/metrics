package agent

import (
	"flag"
	"github.com/RomanenkoDR/metrics/internal/config/agent/agn_types"
	"github.com/caarlos0/env"
	"strings"
)

func ParseOptions() (agn_types.OptionsAgent, error) {
	var opt agn_types.OptionsAgent
	opt.Encrypt = false

	// Чтение параметра командной строки для интервала сбора метрик (по умолчанию 2 секунды)
	flag.IntVar(&opt.PollInterval,
		"p",
		2,
		"Frequency in seconds for collecting metrics")

	// Чтение параметра командной строки для интервала отправки метрик (по умолчанию 10 секунд)
	flag.IntVar(&opt.ReportInterval,
		"r",
		10,
		"Frequency in seconds for sending report to the server")

	// Чтение параметра командной строки для адреса сервера (по умолчанию "localhost:8080")
	flag.StringVar(&opt.ServerAddress,
		"a",
		"localhost:8080",
		"Address of the server to send metrics")

	// Чтение параметра командной строки для установки JWT токена
	flag.StringVar(&opt.Key,
		"k",
		"",
		"Token auth by JWT")

	flag.IntVar(&opt.RateLimit,
		"l",
		3,
		"Rate Limit")

	// Парсинг аргументов командной строки
	flag.Parse()

	opt.ServerAddress = strings.Join([]string{"http:/", opt.ServerAddress, "updates/"}, "/")

	if opt.Key != "" {
		opt.Encrypt = true
		opt.KeyByte = []byte(opt.Key)
	}

	// Парсинг переменных окружения и их присвоение в структуру Options
	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	// Возвращаем структуру с параметрами и nil (ошибки нет)
	return opt, nil
}
