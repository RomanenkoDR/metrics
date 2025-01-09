package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type Database struct {
	Conn *pgx.Conn
}

// Connect устанавливает подключение к базе данных.
func Connect(connstring string) (*Database, error) {
	ctx := context.Background()
	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully")
	return &Database{Conn: conn}, nil // Возвращаем указатель
}

// Close завершает соединение с базой данных.
func (db *Database) Close() {
	db.Conn.Close(context.Background())
}
