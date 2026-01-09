package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	shutdownTimeout := 30 * time.Second
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if parsed, err := time.ParseDuration(timeout); err == nil {
			shutdownTimeout = parsed
		}
	}

	return &Config{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
	}
}
