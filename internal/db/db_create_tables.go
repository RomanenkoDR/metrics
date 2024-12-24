package db

import (
	"context"
)

func (db *Database) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS gauge_metrics (
			id serial PRIMARY KEY,
			name text NOT NULL,
			value double precision NOT NULL,
			timestamp timestamp DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS counter_metrics (
			id serial PRIMARY KEY,
			name text NOT NULL,
			value integer NOT NULL,
			timestamp timestamp DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := db.Conn.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}

	return nil
}
