package server

import (
	"context"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Run запускает сервер с обработкой метрик и поддержкой профилирования.
func Run() {
	runPprof() // Запуск профилировщика pprof

	logger.Info("Запуск сервера...") // Логируем старт сервера

	var store storage.StorageWriter // Хранилище данных

	cfg, err := parseOptions() // Парсим параметры
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprint("флаг на сервере: ", cfg.Key))

	log.Println("Параметры конфигурации сервера: ", zap.Any("metrics", cfg))

	h := handlers.NewHandler() // Создаём новый обработчик запросов

	if cfg.DBDSN != "" { // Проверяем подключение к базе данных
		log.Println("Подключение к базе данных DSN:", cfg.DBDSN)
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Fatalf("Ошибка подключения к базе данных: %v", err)
		}
		logger.Info("Успешное подключение к базе данных")
		store = &database
		h.DBconn = database.Conn
	} else {
		store = &storage.Localfile{Path: cfg.Filename}
	}

	router, err := routers.InitRouter(cfg, h) // Инициализация маршрутизатора
	if err != nil {
		panic(err)
	}

	if cfg.Restore { // Восстановление данных
		err := store.RestoreData(&h.Store)
		if err != nil {
			log.Println("Не удалось восстановить данные: ", err)
		}
	}

	go func() { // Периодическое сохранение данных
		for {
			store.Save(cfg.Interval, h.Store)
		}
	}()

	server := http.Server{ // Определение параметров сервера
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Println("Входящие запросы по: ", cfg.Address)

	idleConnectionsClosed := make(chan struct{}) // Настройка завершения работы сервера
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Остановка сервера")
		if err := store.Write(h.Store); err != nil {
			log.Printf("Ошибка сохранения данных: %v", err)
		}
		defer store.Close()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Ошибка завершения работы HTTP сервера: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	log.Fatal(server.ListenAndServe())
	<-idleConnectionsClosed
	log.Println("Сервер остановлен")
}
