package main

import (
	"log"
	"net/http"

	serverConfigPcg "github.com/RomanenkoDR/metrics/internal/config/serverConfig"
	"github.com/RomanenkoDR/metrics/internal/routers"
)

func main() {
	// Парсинг опций командной строки
	cfg := serverConfigPcg.ParseOptions()

	log.Println("Запуск сервера...")

	// Запуск сервера
	log.Fatal(http.ListenAndServe(cfg.Address, routers.InitRouter()))
}
