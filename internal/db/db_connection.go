package db

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"github.com/jackc/pgx/v5"
	"log"
)

// Connect устанавливает подключение к базе данных и проверяет наличие таблиц.
func Connect(connstring string) (*db_types.Database, error) {
	ctx := context.Background()
	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	db := &db_types.Database{Conn: conn}
	log.Println("Connected to the database successfully")

	// Проверка и создание таблиц при необходимости
	err = createTables(db)
	if err != nil {
		log.Printf("Ошибка при создании таблиц: %v\n", err)
		return nil, err
	}

	log.Println("Все таблицы успешно проверены и/или созданы")
	return db, nil
}

// Close завершает соединение с базой данных.
func Close(db *db_types.Database) {
	db.Conn.Close(context.Background())
}
