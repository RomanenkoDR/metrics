package db

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"github.com/RomanenkoDR/metrics/internal/storage"
)

// RestoreData восстанавливает данные из базы данных в память.
func RestoreData(db *db_types.Database, s *storage.MemStorage) error {
	ctx := context.Background()

	// Восстановление gauge данных
	rows, err := db.Conn.Query(ctx, `SELECT name, value FROM gauge_metrics`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value storage.Gauge
		if err := rows.Scan(&name, &value); err != nil {
			return err
		}
		s.UpdateGauge(name, value)
	}

	// Восстановление counter данных
	rows, err = db.Conn.Query(ctx, `SELECT name, value FROM counter_metrics`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value storage.Counter
		if err := rows.Scan(&name, &value); err != nil {
			return err
		}
		s.UpdateCounter(name, value)
	}

	return nil
}
