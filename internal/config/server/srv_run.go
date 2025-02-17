package server

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/crypto"
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

func Run() {
	// Проверяем и генерируем ключи, если их нет
	if err := crypto.GenerateKeys(); err != nil {
		logger.Fatal("Ошибка генерации ключей", zap.Error(err))
	}

	// Запуск профилировщика pprof
	runPprof()

	// Логируем старт сервера
	logger.Info("Запуск сервера...")

	// Объявляем переменную для хранилища данных
	var store storage.StorageWriter

	// Парсим параметры командной строки и конфигурацию сервера
	cfg, err := parseOptions()
	if err != nil {
		panic(err)
	}

	// Логируем полученные параметры конфигурации
	log.Println("Параметры конфигурации сервера: ", zap.Any("metrics", cfg))

	// Создаём новый обработчик запросов (handler), который будет управлять маршрутами и логикой обработки
	h := handlers.NewHandler()

	// Если в конфигурации указан DSN для подключения к базе данных, то подключаемся к базе
	log.Println("Подключения к базе данных DBDSN сервера:", cfg.DBDSN)
	if cfg.DBDSN != "" {
		// Логируем процесс подключения к базе данных
		log.Println("Подключение к базе данных DSN:", cfg.DBDSN)
		database, err := db.Connect(cfg.DBDSN)
		if err != nil {
			log.Fatalf("Ошибка подключения к базе данных: %v", err)
		} else {
			logger.Info("Успешное подключение к базе данных")
		}

		// Устанавливаем базу данных в качестве хранилища данных
		store = &database

		// Передаём подключение к базе данных в обработчик запросов
		h.DBconn = database.Conn

	} else {
		// Если DSN для базы данных не указан, используем файл для хранения метрик
		store = &storage.Localfile{Path: cfg.Filename}
	}

	// Инициализируем маршрутизатор с конфигурацией и хэндлером

	router, err := routers.InitRouter(cfg, h)
	if err != nil {
		panic(err)
	}

	// Если в конфигурации указан флаг "Restore", восстанавливаем данные из хранилища (файла или базы данных)
	if cfg.Restore {
		err := store.RestoreData(&h.Store)
		// Логируем ошибку восстановления данных, если она произошла
		if err != nil {
			log.Println("Не удалось восстановить данные: ", err)
		}
	}

	// Запускаем горутину для периодической записи данных в хранилище (файл или БД).
	// Интервал указывается в конфигурации.
	go func() {
		for {
			store.Save(cfg.Interval, h.Store)
		}
	}()

	// Определяем параметры HTTP-сервера
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Логируем, что сервер начал слушать входящие запросы на указанном адресе
	log.Println("Входящие запросы по: ", cfg.Address)
	log.Println("Запуск сервера ")

	// Настройка корректного завершения работы сервера
	idleConnectionsClosed := make(chan struct{}) // Канал для оповещения о закрытии всех соединений
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint // Ожидаем поступления сигнала
		// Логируем начало процесса завершения работы
		log.Println("Остановка сервера")

		// Сохраняем оставшиеся данные перед завершением работы
		if err := store.Write(h.Store); err != nil {
			// Логируем ошибку, если не удалось сохранить данные
			log.Printf("Ошибка сохранения даннных: %v", err)
		}

		// Закрываем хранилище (файл или БД)
		defer store.Close()

		// Завершаем работу HTTP-сервера
		if err := server.Shutdown(context.Background()); err != nil {
			// Логируем ошибку завершения сервера, если она произошла
			log.Printf("Ошибка завершения работы HTTP сервера: %v", err)
		}
		// Оповещаем, что все соединения закрыты
		close(idleConnectionsClosed)
	}()

	// Запускаем сервер для прослушивания входящих запросов
	log.Fatal(server.ListenAndServe())

	// Ожидаем закрытия всех соединений перед завершением программы
	<-idleConnectionsClosed
	// Логируем завершение работы сервера
	log.Println("Сервер остановлен")
}
