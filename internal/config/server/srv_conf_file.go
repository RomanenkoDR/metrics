package server

import (
	"encoding/json"
	"os"
	"time"
)

// serverFileConfig хранит конфигурацию сервера из JSON-файла
type serverFileConfig struct {
	Address       string        `json:"address"`        // Адрес сервера
	Restore       bool          `json:"restore"`        // Восстанавливать ли метрики из файла
	StoreInterval time.Duration `json:"store_interval"` // Интервал сохранения метрик
	StoreFile     string        `json:"store_file"`     // Файл хранения метрик
	DatabaseDSN   string        `json:"database_dsn"`   // Строка подключения к БД
	CryptoKey     string        `json:"crypto_key"`     // Путь к приватному ключу
}

// loadConfigFromFile загружает конфигурацию сервера из JSON-файла
func loadConfigFromFile(path string) (*serverFileConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &serverFileConfig{}
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
