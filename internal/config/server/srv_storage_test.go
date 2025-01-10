package server

import (
	"github.com/RomanenkoDR/metrics/internal/config/server/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupStorage_FileStorage(t *testing.T) {
	cfg := types.OptionsServer{
		Filename: "/tmp/test_metrics.json",
		DBDSN:    "",
	}

	store, err := setupStorage(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, store, "Хранилище не должно быть nil")
}

func TestSetupStorage_Database(t *testing.T) {
	cfg := types.OptionsServer{
		Filename: "/tmp/test_metrics.json",
		DBDSN:    "user=postgres dbname=test sslmode=disable",
	}

	store, err := setupStorage(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, store, "Хранилище не должно быть nil")
}
