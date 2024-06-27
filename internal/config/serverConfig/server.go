package serverConfig

import (
	"flag"
	"fmt"
	"os"
)

// Содержит конфигурацию сервера, такую как адрес.
type ServerConfig struct {
	Address string // Адрес сервера
}

// Создает и возвращает новый экземпляр ServerConfig с инициализированными значениями по умолчанию
func NewServerConfiguration() *ServerConfig {
	return &ServerConfig{}
}

// Инициализирует конфигурацию сервера, используя флаги командной строки и переменные окружения.
func (c *ServerConfig) InitServerConfiguration() *ServerConfig {
	// Изменение конфигурации на основе переменной окружения
	if addr := os.Getenv("ADDRESS"); addr != "" {
		c.Address = addr
	}
	// Установка флага командной строки для адреса сервера с значением по умолчанию "localhost:8080"
	flag.StringVar(&c.Address, "a", "localhost:8080", "Адрес и порт сервера для запуска сервера")
	// Парсинг флагов командной строки
	flag.Parse()

	// Проверка на неизвестные флаги
	flag.VisitAll(func(f *flag.Flag) {
		// Если флаг не был распознан (не спарсился), выводим сообщение об ошибке
		if !flag.Parsed() {
			fmt.Printf("неизвестный флаг: %s\n", f.Name)
			// Выводим информацию о правильном использовании флагов
			flag.Usage()
		}
	})
	// Возвращаем указатель на инициализированный обьект ServerConfig
	return c

}
