package main

import (
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/server"
	"github.com/RomanenkoDR/metrics/internal/crypto"
	_ "net/http/pprof"
	"os"
)

var (
	buildVersion = "N/A" // VERSION 1-курс, 8-спринт, 20-инкремент
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// Запуск сборки приложения из корня проекта
// go build -ldflags "-X main.buildVersion=1.8.20 -X main.buildDate=$(date -u +%Y-%m-%d) -X main.buildCommit=$(git rev-parse HEAD)" -o server cmd/server/main.go
func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// Путь к ключу
	privateKeyPath := "private.pem"
	publicKeyPath := "../agent/public.pem"

	// Проверяем, существуют ли уже ключи
	privateExists := fileExists(privateKeyPath)
	publicExists := fileExists(publicKeyPath)

	if !privateExists || !publicExists {
		fmt.Println("Один из ключей не найден. Генерируем новые...")

		err := crypto.GenerateRSAKeys(2048)
		if err != nil {
			fmt.Println("Ошибка генерации ключей:", err)
			os.Exit(1)
		}
		fmt.Println("Ключи успешно сгенерированы.")
	} else {
		fmt.Println("Оба ключа найдены, генерация не требуется.")
	}

	server.Run()
}

// fileExists проверяет существование файла по указанному пути
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
