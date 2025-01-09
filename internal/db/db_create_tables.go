package db

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"log"
)

// createTables создает таблицы в базе данных.
func createTables(db *db_types.Database) error {
	tables := []db_types.TableConfig{
		{
			Name: "gauge_metrics",
			CreateQuery: `CREATE TABLE IF NOT EXISTS gauge_metrics (
				id serial PRIMARY KEY,
				name text UNIQUE,
				value double precision,
				timestamp timestamp
			)`,
		},
		{
			Name: "counter_metrics",
			CreateQuery: `CREATE TABLE IF NOT EXISTS counter_metrics (
				id serial PRIMARY KEY,
				name text UNIQUE,
				value integer,
				timestamp timestamp
			)`,
		},
	}

	for _, table := range tables {
		_, err := db.Conn.Exec(context.Background(), table.CreateQuery)
		if err != nil {
			log.Printf("Ошибка создания таблицы '%s': %v\n", table.Name, err)
			return err
		}
	}

	log.Println("Все таблицы успешно созданы или уже существуют")
	return nil
}
