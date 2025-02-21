package agent

import (
	"encoding/json"
	"os"
	"time"
)

// agentConfFromFile хранит конфигурацию агента из JSON-файла
type agentConfFromFile struct {
	ServerAddress  string        `json:"address"`         // Адрес сервера для отправки метрик
	ReportInterval time.Duration `json:"report_interval"` // Интервал отправки метрик на сервер
	PollInterval   time.Duration `json:"poll_interval"`   // Интервал сбора метрик
	CryptoKey      string        `json:"crypto_key"`      // Путь к публичному ключу для шифрования
}

// loadConfigFromFile загружает конфигурацию агента из JSON-файла
func loadConfigFromFile(path string) (*agentConfFromFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &agentConfFromFile{}
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
