package db

import "github.com/RomanenkoDR/metrics/internal/db/types"

// TableConfigs возвращает конфигурации таблиц для базы данных.
func tableConfigs() []types.TableConfig {
	return []types.TableConfig{
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
}
