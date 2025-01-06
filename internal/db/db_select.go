package db

import (
	"context"
	"fmt"
)

func (db *Database) SelectAll() error {
	rows, err := db.Conn.Query(context.Background(),
		`SELECT * FROM counter_metrics
                           UNION
                           SELECT * FROM gauge_metrics
                           Limit 10`)
	if err != nil {
		return err
	}

	fmt.Println("SelectAll")
	for rows.Next() {
		data, err := rows.Values()
		if err != nil {
			return err
		}
		fmt.Println(data)
	}

	return nil
}
