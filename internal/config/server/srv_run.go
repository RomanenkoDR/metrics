package server

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	runPprof()

	// Объявляем переменную для хранилища данных
	var store storage.StorageWriter

	logger.Info("Запуск сервера...")

	// Парсим параметры командной строки
	cfg, err := parseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов", zap.Error(err))
	}

	// Создаём новый обработчик запросов
	h := handlers.NewHandler()

	// Проверяем, указан ли приватный ключ
	if cfg.CryptoKey == "" {
		logger.Warn("Флаг -crypto-key не задан, сервер не сможет расшифровывать данные")
	} else {
		logger.Info("Используется приватный ключ для расшифровки", zap.String("cryptoKey", cfg.CryptoKey))
		h.SetCryptoKey(cfg.CryptoKey)
	}

	// Определяем хранилище данных (БД или файл)
	if cfg.DBDSN != "" {
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
		}
		logger.Info("Успешное подключение к базе данных")
		store = &database
		h.DBconn = database.Conn
	} else {
		store = &storage.Localfile{Path: cfg.Filename}
	}

	// Создаём новое хранилище данных
	h.Store = storage.New()

	// Восстанавливаем данные, если они есть
	err = store.RestoreData(&h.Store)
	if err != nil {
		logger.Warn("Не удалось восстановить данные из хранилища", zap.Error(err))
	} else {
		logger.Info("Данные успешно загружены из хранилища")
	}

	// Инициализируем маршрутизатор
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		logger.Fatal("Ошибка инициализации маршрутизатора", zap.Error(err))
	}

	// Запускаем сервер
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Запускаем периодическое сохранение данных
	go func() {
		for {
			store.Save(cfg.Interval, h.Store)
		}
	}()

	// Обрабатываем сигналы завершения работы сервера
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		logger.Info("Остановка сервера")

		// Сохраняем данные перед выходом
		if err := store.Write(h.Store); err != nil {
			logger.Error("Ошибка сохранения данных перед выходом", zap.Error(err))
		}

		store.Close()

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Ошибка завершения сервера", zap.Error(err))
		}
		close(idleConnectionsClosed)
	}()

	logger.Info("Сервер запущен", zap.String("address", cfg.Address))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}

	<-idleConnectionsClosed
	logger.Info("Сервер остановлен")
}
