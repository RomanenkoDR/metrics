package serverConfig

import (
	"flag"
	"os"
)

// Options содержит параметры конфигурации для сервера:
// адрес, на котором сервер будет слушать входящие соединения.
type Options struct {
	Address string
}

// ParseOptions парсит опции из командной строки и переменных окружения
func ParseOptions() Options {
	var cfg Options

	// Установка значения по умолчанию и описания для опции командной строки
	flag.StringVar(&cfg.Address,
		"a", "localhost:8080",
		"Укажите адрес и порт в формате <address>:<port>")
	flag.Parse()

	// Получение переменных окружения
	// Проверяет наличие переменной окружения ADDRESS.
	// Если переменная установлена, её значение используется для адреса сервера.
	if a := os.Getenv("ADDRESS"); a != "" {
		cfg.Address = a
	}

	return cfg
}
