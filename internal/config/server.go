package config

import (
	"flag"
	"os"
)

func getEnvOrDefault(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

type ServerConfig struct {
	Address string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: getEnvOrDefault("ADDRESS", "localhost:8080"),
	}
}

func (c *ServerConfig) Init() {
	flag.StringVar(&c.Address, "a", c.Address, "address and port to run server")
	flag.Parse()
}
