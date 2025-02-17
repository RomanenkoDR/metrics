package main

import (
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// go build -ldflags "-X main.buildVersion=1.8.20 -X main.buildDate=$(date -u +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o server cmd/server/main.go
func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// Запускаем сервер
	server.Run()

}
