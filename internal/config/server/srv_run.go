package server

import (
	"context"
	//"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	// Логируем старт сервера
	logger.Info("Запуск сервера...")

	// Парсим параметры командной строки и конфигурацию сервера
	cfg, err := parseOptions()
	if err != nil {
		logger.Fatal("Ошибка разбора флагов", zap.Error(err))
	}

	// Создаём новый обработчик запросов
	h := handlers.NewHandler()

	// Проверяем переданный путь к приватному ключу
	if cfg.CryptoKey == "" {
		logger.Warn("Флаг -crypto-key не задан, сервер не сможет расшифровывать данные")
	} else {
		logger.Info("Используется приватный ключ для расшифровки", zap.String("cryptoKey", cfg.CryptoKey))
		h.SetCryptoKey(cfg.CryptoKey)
	}

	// Подключаемся к базе данных, если указан DSN
	if cfg.DBDSN != "" {
		dbConn, err := db.Connect(cfg.DBDSN)
		if err != nil {
			logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
		}
		logger.Info("Успешное подключение к базе данных")
		h.DBconn = dbConn.Conn
	}

	// Инициализируем маршрутизатор
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		logger.Fatal("Ошибка инициализации маршрутизатора", zap.Error(err))
	}

	// Определяем HTTP-сервер
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Обрабатываем сигналы завершения работы сервера
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint // Ждём сигнал завершения

		logger.Info("Остановка сервера")

		// Закрываем соединение с базой данных
		if h.DBconn != nil {
			h.DBconn.Close(context.Background())
		}

		// Завершаем работу HTTP-сервера
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Ошибка завершения сервера", zap.Error(err))
		}
		close(idleConnectionsClosed)
	}()

	// Запускаем сервер
	logger.Info("Сервер запущен", zap.String("address", cfg.Address))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}

	<-idleConnectionsClosed
	logger.Info("Сервер остановлен")
}
