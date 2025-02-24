package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanenkoDR/metrics/internal/db"
	"github.com/RomanenkoDR/metrics/internal/handlers"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/routers"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"go.uber.org/zap"
)

func Run() {
	runPprof()

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
	var store storage.StorageWriter

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
	if err := store.RestoreData(&h.Store); err != nil {
		logger.Warn("Не удалось восстановить данные из хранилища", zap.Error(err))
	} else {
		logger.Info("Данные успешно загружены из хранилища")
	}

	// Инициализируем маршрутизатор
	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		logger.Fatal("Ошибка инициализации маршрутизатора", zap.Error(err))
	}

	// Создаём HTTP-сервер
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Контекст с отменой для graceful shutdown
	// Контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Гарантируем вызов cancel() перед выходом

	// Канал для перехвата сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Канал завершения работы сервера
	done := make(chan struct{})

	// Запускаем сервер в горутине
	go func() {
		logger.Info("Сервер запущен", zap.String("address", cfg.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	// Запускаем периодическое сохранение данных в отдельной горутине
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("Остановка фонового сохранения данных")
				return
			case <-ticker.C:
				logger.Debug("Автосохранение данных")
				store.Save(cfg.Interval, h.Store)
			}
		}
	}()

	// Ожидаем сигнала завершения
	go func() {
		sig := <-sigChan
		logger.Info("Получен сигнал завершения", zap.String("signal", sig.String()))
		cancel() // Теперь вызываем cancel() при завершении работы сервера

		// Создаём контекст с таймаутом для shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		// Завершаем HTTP-сервер
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Ошибка завершения сервера", zap.Error(err))
		}

		// Сохраняем все несохранённые данные
		logger.Info("Сохранение данных перед выходом")
		if err := store.Write(h.Store); err != nil {
			logger.Error("Ошибка сохранения данных", zap.Error(err))
		}

		// Закрываем хранилище, если оно использует БД
		store.Close()

		close(done)
	}()

	// Ожидаем завершения работы
	<-done
	logger.Info("Сервер остановлен")

}
