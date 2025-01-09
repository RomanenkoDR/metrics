package types

// TableConfig описывает конфигурацию для создания таблицы.
type TableConfig struct {
	Name        string // Название таблицы
	CreateQuery string // SQL-запрос для создания таблицы
}
