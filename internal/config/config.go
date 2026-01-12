package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
	MongoURI        string
	MongoDB         string
	RedisAddr       string
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

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "godemo"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return &Config{
		Port:            port,
		ShutdownTimeout: shutdownTimeout,
		MongoURI:        mongoURI,
		MongoDB:         mongoDB,
		RedisAddr:       redisAddr,
	}
}
