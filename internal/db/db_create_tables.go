package db

import (
	"context"
	"log"
)

func (db *Database) createTables() error {
	_, err := db.Conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS gauge_metrics(
        id serial PRIMARY KEY,
        name text,
        value double precision,
        timestamp timestamp)`)
	if err != nil {
		log.Println("Error creating gauge_metrics table:", err)
		return err
	}

	_, err = db.Conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS counter_metrics(
        id serial PRIMARY KEY,
        name text,
        value integer,
        timestamp timestamp)`)
	if err != nil {
		log.Println("Error creating counter_metrics table:", err)
		return err
	}

	log.Println("Tables created successfully")
	return nil
}
