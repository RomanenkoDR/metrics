package db

import (
	"context"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"log"
)

// SelectAll выполняет выборку данных из таблиц базы данных.
func SelectAll(db *db_types.Database) error {
	query := `
		SELECT name, value, timestamp FROM counter_metrics
		UNION ALL
		SELECT name, value, timestamp FROM gauge_metrics
		LIMIT 10
	`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Ошибка выполнения SelectAll: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value interface{}
		var timestamp interface{}
		if err := rows.Scan(&name, &value, &timestamp); err != nil {
			log.Printf("Ошибка обработки строки: %v", err)
			return err
		}
		fmt.Printf("Name: %s, Value: %v, Timestamp: %v\n", name, value, timestamp)
	}

	return nil
}
