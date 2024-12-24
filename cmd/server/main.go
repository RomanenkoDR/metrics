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
	// Store variable will be used file or database to save metrics
	var store storage.StorageWriter

	// Parse cli options into config
	cfg, err := server.ParseOptions()
	if err != nil {
		panic(err)
	}

	log.Println("Params:", cfg)

	// Handler for router
	h := handlers.NewHandler()

	// Identify wether use DB or file to save metrics
	if cfg.DBDSN != "" {
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Println(err)
		}

		// Use database as a store
		store = &database

		//Define DB for handlers
		h.DBconn = database.Conn

	} else {
		// use json file to store metrics
		store = &storage.Localfile{Path: cfg.Filename}
	}

	// Init router
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		panic(err)
	}

	if cfg.Restore {
		err := store.RestoreData(h.Store)
		if err != nil {
			log.Println("Could not restore data: ", err)
		}
	}

	// Write MemStorage to a store provider
	// Interval used for file saving
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

	log.Println("Started. Running")

	// Graceful shutdown
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server")

		if err := store.Write(h.Store); err != nil {
			log.Printf("Error during saving data to file: %v", err)
		}

		// Close file/db
		defer store.Close()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	// Run server
	log.Fatal(server.ListenAndServe())

	<-idleConnectionsClosed
	log.Println("Server shutdown")
}
