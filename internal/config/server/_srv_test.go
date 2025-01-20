package server

import (
	"flag"
	"os"
	"testing"

	"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/stretchr/testify/assert"
)

// TestParseOptions проверяет корректность парсинга конфигурации.
func TestParseOptions(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("ADDRESS", "127.0.0.1:9000")
	os.Setenv("INTERVAL", "500")
	os.Setenv("FILENAME", "/tmp/test_metrics.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("KEY", "test-key")
	os.Setenv("DBDSN", "postgres://user:pass@localhost/dbname")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("INTERVAL")
		os.Unsetenv("FILENAME")
		os.Unsetenv("RESTORE")
		os.Unsetenv("KEY")
		os.Unsetenv("DBDSN")
	}()

	// Сброс флагов перед тестом (необходим для тестов flag)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cfg, err := ParseOptions()
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:9000", cfg.Address)
	assert.Equal(t, 500, cfg.Interval)
	assert.Equal(t, "/tmp/test_metrics.json", cfg.Filename)
	assert.Equal(t, false, cfg.Restore)
	assert.Equal(t, "test-key", cfg.Key)
	assert.Equal(t, "postgres://user:pass@localhost/dbname", cfg.DBDSN)
}

// TestSetupStorage проверяет настройку хранилища.
func TestSetupStorage(t *testing.T) {
	t.Run("Database storage", func(t *testing.T) {
		cfg := types.OptionsServer{
			DBDSN: "postgres://user:pass@localhost/dbname",
		}
		storage, err := setupStorage(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})

	t.Run("File storage", func(t *testing.T) {
		cfg := types.OptionsServer{
			Filename: "/tmp/test_metrics.json",
		}
		storage, err := setupStorage(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})
}

// TestSetupShutdown тестирует корректное завершение работы сервера.
func TestSetupShutdown(t *testing.T) {
	// Тестирование setupShutdown возможно только через интеграционные тесты с имитацией сигналов.
	// Для этого лучше использовать тестовые библиотеки, такие как "os/signal" с mock.
	t.Skip("Тестирование graceful shutdown пропущено, требуется интеграционное окружение")
}
