package server

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseOptions(t *testing.T) {
	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("STORE_INTERVAL", "10")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/metrics.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("DATABASE_DSN", "user=postgres dbname=test sslmode=disable")
	os.Setenv("KEY", "testkey")

	cfg, err := ParseOptions()

	assert.NoError(t, err, "ParseOptions должна возвращать nil-ошибку")
	assert.Equal(t, "localhost:9090", cfg.Address)
	assert.Equal(t, 10, cfg.Interval)
	assert.Equal(t, "/tmp/metrics.json", cfg.Filename)
	assert.Equal(t, false, cfg.Restore)
	assert.Equal(t, "user=postgres dbname=test sslmode=disable", cfg.DBDSN)
	assert.Equal(t, "testkey", cfg.Key)
}
