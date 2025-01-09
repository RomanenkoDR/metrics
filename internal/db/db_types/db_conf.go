package db_types

import "github.com/jackc/pgx/v5"

// Database представляет структуру для подключения к базе данных.
type Database struct {
	Conn *pgx.Conn
}

// TableConfig описывает конфигурацию для создания таблицы.
type TableConfig struct {
	Name        string
	CreateQuery string
}
