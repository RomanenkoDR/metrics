package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type Database struct {
	Conn *pgx.Conn
}

func Connect(connstring string) (*Database, error) {
	var db Database

	ctx := context.Background()
	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		return nil, err
	}

	db.Conn, err = pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully")

	err = db.createTables()
	if err != nil {
		return nil, err
	}

	log.Println("Tables created or verified successfully")
	return &db, nil
}

func (db *Database) Close() {
	db.Conn.Close(context.Background())
}
