package main

import "github.com/RomanenkoDR/metrics/internal/config/server"
import _ "net/http/pprof"

// main - основная функция для запуска агента.
// Запускает функцию server.Run, которая отвечает за выполнение всех операций агента.
func main() {
	server.Run()
}
