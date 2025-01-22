package server

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func runPprof() { // В функции Run
	go func() {
		log.Println("pprof запущен на :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Не удалось запустить pprof: %v", err)
		}
	}()
}
