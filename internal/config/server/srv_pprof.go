package server

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// runPprof запускает встроенный профилировщик pprof на порту 6060.
func runPprof() {
	go func() {
		log.Println("pprof запущен на :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Не удалось запустить pprof: %v", err)
		}
	}()
}
