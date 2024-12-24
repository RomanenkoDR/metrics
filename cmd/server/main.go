package main

import (
	"context"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/server"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting server...")
	var store storage.StorageWriter

	cfg, err := server.ParseOptionsServer()
	if err != nil {
		log.Fatalf("Failed to parse server options: %v", err)
	}

	// Установить порт из переменной окружения, если указан
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		cfg.Address = ":" + port
	}

	log.Println("Params:", cfg)

	h := handlers.NewHandler()

	// Подключение к базе данных или локальному хранилищу
	if cfg.DBDSN != "" {
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Println("Error connecting to database, switching to local file storage:", err)
			store = &storage.Localfile{Path: cfg.Filename}
		} else {
			store = &database
			h.DBconn = database.Conn
		}
	} else {
		store = &storage.Localfile{Path: cfg.Filename}
	}

	// Инициализация маршрутизатора
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	// Восстановление данных
	if cfg.Restore {
		if err := store.RestoreData(h.Store); err != nil {
			log.Println("Could not restore data: ", err)
		}
	}

	// Контекст для управления временем жизни горутин
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Цикл сохранения данных
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Save goroutine exiting...")
				return
			case <-ticker.C:
				if err := store.Save(cfg.Interval, h.Store); err != nil {
					log.Println("Error saving data:", err)
				}
			}
		}
	}(ctx)

	// Инициализация HTTP-сервера
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Проверка готовности сервера
	if err := waitForServerReady(cfg.Address, 5*time.Second); err != nil {
		log.Fatalf("Server not ready: %v", err)
	}

	log.Println("Started. Running")

	idleConnectionsClosed := make(chan struct{})

	// Обработка сигналов завершения
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server")

		if err := store.Write(h.Store); err != nil {
			log.Printf("Error during saving data to file: %v", err)
		}

		store.Close()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	// Запуск сервера
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP Server Listen Error: %v", err)
	}

	<-idleConnectionsClosed
	log.Println("Server shutdown")
}

// Проверка готовности сервера
func waitForServerReady(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not ready on %s within %v", address, timeout)
}
