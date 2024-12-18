package server

import (
	"flag"
	"github.com/caarlos0/env"
)

type Options struct {
	Address  string `env:"ADDRESS"`
	Interval int    `env:"STORE_INTERVAL"`
	Filename string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
	DBDSN    string `env:"DATABASE_DSN"`
	Key      string `env:"KEY"`
}

func ParseOptions() (Options, error) {
	var cfg Options

	// Чтение флага "-a" для задания адреса сервера и порта
	flag.StringVar(&cfg.Address,
		"a",
		"localhost:8080",
		"Add address and port in format <address>:<port>")

	// Чтение флага "-i" для задания интервала сохранения метрик в файл
	flag.IntVar(&cfg.Interval,
		"i",
		300,
		"Saving metrics to file interval")

	// Чтение флага "-f" для задания пути к файлу, где будут храниться метрики
	flag.StringVar(&cfg.Filename,
		"f",
		"./metrics.json",
		"File path")

	// Чтение флага "-r" для задания опции восстановления метрик из файла
	flag.BoolVar(&cfg.Restore,
		"r",
		true,
		"Restore metrics value from file")

	// Чтение флака "-k" для задания токена JWT
	flag.StringVar(&cfg.Key,
		"k",
		"",
		"Token auth by JWT")

	// Чтение флага "-d" для задания строки подключения к базе данных
	flag.StringVar(&cfg.DBDSN,
		"d",
		"",
		"Connection string in Postgres format")

	// Парсинг флагов командной строки
	flag.Parse()

	// Получение значений из переменных окружения и их применение
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
