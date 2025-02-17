package server

import (
   "flag"
   "github.com/RomanenkoDR/metrics/internal/config/server/types"
   "github.com/caarlos0/env"
)

// parseOptions загружает конфигурацию сервера
func parseOptions() (types.Options, error) {
   var cfg types.Options

   // Читаем флаги командной строки
   flag.StringVar(&cfg.Address, "a", "localhost:8080", "Server address and port")
   flag.IntVar(&cfg.Interval, "i", 300, "Interval for saving metrics to file")
   flag.StringVar(&cfg.Filename, "f", "./metrics.json", "Path to metrics file")
   flag.BoolVar(&cfg.Restore, "r", true, "Restore metrics from file")
   flag.StringVar(&cfg.Key, "k", "", "JWT authentication key")
   flag.StringVar(&cfg.DBDSN, "d", "", "Database connection string")
   flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "Path to private key file") // Единственное объявление

   flag.Parse()

   // Загружаем конфиг из переменных окружения
   err := env.Parse(&cfg)
   if err != nil {
	  return cfg, err
   }

   return cfg, nil
}
