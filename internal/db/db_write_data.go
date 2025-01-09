package db

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Write сохраняет или обновляет данные в базе данных.
func Write(db *db_types.Database, s storage.MemStorage) error {
	ctx := context.Background()

	for k, v := range s.CounterData {
		_, err := db.Conn.Exec(ctx,
			`INSERT INTO counter_metrics (name, value, timestamp) 
			 VALUES ($1, $2, $3) 
			 ON CONFLICT (name) DO UPDATE SET value = $2, timestamp = $3`,
			k, v, time.Now())
		if err != nil {
			return err
		}
	}

	for k, v := range s.GaugeData {
		_, err := db.Conn.Exec(ctx,
			`INSERT INTO gauge_metrics (name, value, timestamp) 
			 VALUES ($1, $2, $3) 
			 ON CONFLICT (name) DO UPDATE SET value = $2, timestamp = $3`,
			k, v, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
