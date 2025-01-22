package main

import (
	"github.com/RomanenkoDR/metrics/internal/config/agent"
)

// main - основная функция для запуска агента.
// Запускает функцию agent.Run, которая отвечает за выполнение всех операций агента.
func main() {
	agent.Run()
}
