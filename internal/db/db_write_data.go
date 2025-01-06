package db

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

func (db *Database) Write(s storage.MemStorage) error {
	for k, v := range s.CounterData {
		_, err := db.Conn.Exec(context.Background(),
			`INSERT INTO counter_metrics (name, value, timestamp) VALUES ($1, $2, $3)`,
			k, v, time.Now())
		if err != nil {
			return err
		}
	}

	for k, v := range s.GaugeData {
		_, err := db.Conn.Exec(context.Background(),
			`INSERT INTO gauge_metrics (name, value, timestamp) VALUES ($1, $2, $3)`,
			k, v, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Save(t int, s storage.MemStorage) error {
	time.Sleep(time.Second * time.Duration(t))
	return db.Write(s)
}
