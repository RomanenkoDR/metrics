package server

import (
	"context"
	"fmt"
	"net"
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

	// Проверяем доступность порта перед запуском
	ln, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Fatal("Порт занят, невозможно запустить сервер", zap.String("address", cfg.Address), zap.Error(err))
	}
	ln.Close() // Освобождаем порт перед запуском сервера

	// Создаём новый обработчик запросов
	h := handlers.NewHandler()

	// Устанавливаем приватный ключ для расшифровки (если указан)
	if cfg.CryptoKey != "" {
		logger.Info("Используется приватный ключ для расшифровки", zap.String("cryptoKey", cfg.CryptoKey))
		h.SetCryptoKey(cfg.CryptoKey)
	} else {
		logger.Warn("Флаг -crypto-key не задан, сервер не сможет расшифровывать данные")
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

	// Инициализируем хранилище данных
	h.Store = storage.New()

	// Восстанавливаем данные
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Запускаем фоновое сохранение данных
	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Остановка фонового сохранения данных")
				ticker.Stop()
				return
			case <-ticker.C:
				logger.Debug("Автосохранение данных")
				if err := store.Save(cfg.Interval, h.Store); err != nil {
					logger.Error("Ошибка автосохранения", zap.Error(err))
				}
			}
		}
	}()

	// Ожидаем сигнал завершения
	go func() {
		sig := <-sigChan
		logger.Info("Получен сигнал завершения", zap.String("signal", sig.String()))
		cancel()

		// Контекст с таймаутом для graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		logger.Info("Ожидание завершения сервера...")
		errCh := make(chan error, 1)

		go func() {
			errCh <- server.Shutdown(shutdownCtx)
		}()

		select {
		case err := <-errCh:
			if err != nil {
				logger.Error("Ошибка при завершении сервера", zap.Error(err))
			} else {
				logger.Info("Сервер успешно завершился")
			}
		case <-time.After(5 * time.Second): // Принудительный выход
			logger.Error("Принудительное завершение, сервер завис")
			os.Exit(1)
		}

		// Финальное сохранение данных
		logger.Info("Сохранение данных перед выходом")
		if err := store.Write(h.Store); err != nil {
			logger.Error("Ошибка сохранения данных", zap.Error(err))
		}

		// Закрываем соединение с БД (если используется)
		logger.Info("Закрытие соединения с хранилищем")
		store.Close()
		logger.Info("Соединение с хранилищем закрыто")

		// Выводим сообщение, что сервер остановился (для теста `restart_server`)
		fmt.Println("SERVER STOPPED")
		logger.Info("Сервер остановлен, можно запускать новый процесс")

		close(done)
	}()

	// Ожидаем завершения работы сервера
	<-done
	logger.Info("Сервер полностью остановлен")
}
