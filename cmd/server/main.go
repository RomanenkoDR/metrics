package main

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/server"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting server...")

	var store storage.StorageWriter

	cfg, err := server.ParseOptions()
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandler()

	if cfg.DBDSN != "" {
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Println(err)
		}

		store = &database

		h.DBconn = database.Conn

	} else {
		store = &storage.Localfile{Path: cfg.Filename}
	}

	// Init router
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		panic(err)
	}

	if cfg.Restore {
		err := store.RestoreData(&h.Store)
		if err != nil {
			log.Println("Could not restore data: ", err)
		}
	}

	go func() {
		for {
			store.Save(cfg.Interval, h.Store)
		}
	}()

	// Define server parameters
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Println("Сервер запущен")

	// Graceful shutdown
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Завершение работы сервера")

		if err := store.Write(h.Store); err != nil {
			log.Printf("Ошибка сохранения данных в файл: %v", err)
		}

		// Close file/db
		defer store.Close()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server остановлен с ошибкой: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	// Run server
	log.Fatal(server.ListenAndServe())

	<-idleConnectionsClosed
	log.Println("Сервер остановлен")
}
