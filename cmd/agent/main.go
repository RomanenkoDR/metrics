package main

import (
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/agent"
)

var (
	buildVersion = "N/A" // VERSION 1-курс, 8-спринт, 20-инкремент
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// Запуск сборки приложения из корня проекта
// go build -ldflags "-X main.buildVersion=1.8.20 -X main.buildDate=$(date -u +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o agent cmd/agent/main.go

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	agent.Run()
}
