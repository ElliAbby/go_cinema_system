package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

)

type ServerConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type Config struct {
	Server ServerConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables and defaults")
	}

	return &Config{
		Server: ServerConfig{
			Addr:            getEnv("APP_ADDR", ":8080"),
			ReadTimeout:     getDurationEnv("APP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("APP_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("APP_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
	}
}

func getEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}

	return value
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Недопустимая длительность в %s=%q, используется значение по умолчанию %s", key, value, defaultVal)
		return defaultVal
	}

	return duration
}
