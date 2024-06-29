package main

import (
	"log"
	"net/http"

	serverConfigPcg "github.com/RomanenkoDR/metrics/internal/config/serverConfig"
	"github.com/RomanenkoDR/metrics/internal/routers"
)

func main() {
	//parse cli options
	cfg := serverConfigPcg.ParseOptions()

	log.Println("Starting server...")

	//run server
	log.Fatal(http.ListenAndServe(cfg.Address, routers.InitRouter()))
}
