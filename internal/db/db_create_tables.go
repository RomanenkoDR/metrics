package db

import (
	"context"
	"log"
)

// CreateTables проверяет и создает таблицы на основе конфигураций.
func (db *Database) createTables() error {
	configs := tableConfigs()

	for _, cfg := range configs {
		if _, err := db.Conn.Exec(context.Background(), cfg.CreateQuery); err != nil {
			log.Printf("Ошибка создания таблицы '%s': %v\n", cfg.Name, err)
			return err
		}
		log.Printf("Таблица '%s' создана или уже существует\n", cfg.Name)
	}

	return nil
}
