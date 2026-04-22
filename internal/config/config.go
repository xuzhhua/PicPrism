package config

import (
	"os"
)

type Config struct {
	Token   string
	Port    string
	DataDir string
}

func Load() *Config {
	token := os.Getenv("PICPRISM_TOKEN")
	port := os.Getenv("PICPRISM_PORT")
	if port == "" {
		port = "8080"
	}
	dataDir := os.Getenv("PICPRISM_DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	return &Config{
		Token:   token,
		Port:    port,
		DataDir: dataDir,
	}
}
