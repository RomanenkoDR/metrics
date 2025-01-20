package db

import (
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testConnStr = "postgres://user:pass@localhost/dbname" // замените на вашу строку подключения
)

// TestConnect проверяет подключение к базе данных
func TestConnect(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()
}

// TestCreateTables проверяет создание таблиц
func TestCreateTables(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// Попробуем создать таблицы
	err = db.createTables()
	assert.NoError(t, err)
}

// TestWrite проверяет запись данных в базу данных
func TestWrite(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// Создаём пример данных для записи
	memStorage := storage.MemStorage{
		CounterData: map[string]storage.Counter{
			"metric1": 10,
			"metric2": 20,
		},
		GaugeData: map[string]storage.Gauge{
			"metricA": 15.5,
			"metricB": 30.3,
		},
	}

	// Записываем данные
	err = db.Write(memStorage)
	assert.NoError(t, err)
}

// TestSelectAll проверяет выборку данных из базы данных
func TestSelectAll(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// Попробуем выполнить выборку
	err = db.SelectAll()
	assert.NoError(t, err)
}

// TestRestoreData проверяет восстановление данных (метод-заглушка, если нет реализации)
func TestRestoreData(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// Создаём пустой MemStorage для теста
	memStorage := storage.MemStorage{}

	// Проверяем, что метод не вызывает ошибок
	err = db.RestoreData(&memStorage)
	assert.NoError(t, err)
}

// TestSave проверяет метод сохранения данных с задержкой
func TestSave(t *testing.T) {
	db, err := Connect(testConnStr)
	assert.NoError(t, err)
	defer db.Close()

	// Создаём пример данных для записи
	memStorage := storage.MemStorage{
		CounterData: map[string]storage.Counter{
			"metric1": 10,
			"metric2": 20,
		},
		GaugeData: map[string]storage.Gauge{
			"metricA": 15.5,
			"metricB": 30.3,
		},
	}

	// Сохраняем данные с задержкой
	err = db.Save(1, memStorage)
	assert.NoError(t, err)
}
