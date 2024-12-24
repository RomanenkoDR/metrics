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
	"time"
)

func main() {
	log.Println("Starting server...")

	var store storage.WriterStorage

	cfg, err := server.ParseOptions()
	if err != nil {
		panic(err)
	}

	log.Printf("Params: %+v", cfg)

	h := handlers.NewHandler()

	if cfg.DBDSN != "" {
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		store = &database
		h.DBconn = database.Conn
	} else {
		store = &storage.Localfile{Path: cfg.Filename}
	}

	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		panic(err)
	}

	if cfg.Restore {
		log.Println("Restoring metrics from storage...")
		err := store.RestoreData(&h.Store)
		if err != nil {
			log.Printf("Could not restore data: %v", err)
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(cfg.Interval))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("Saving metrics to storage...")
				if err := store.Save(cfg.Interval, h.Store); err != nil {
					log.Printf("Error saving metrics: %v", err)
				}
			}
		}
	}()

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Printf("Started server on %s", cfg.Address)

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server gracefully...")
		if err := store.Write(h.Store); err != nil {
			log.Printf("Error during saving data: %v", err)
		}
		store.Close()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	log.Fatal(server.ListenAndServe())
	<-idleConnectionsClosed
	log.Println("Server shutdown")
}
