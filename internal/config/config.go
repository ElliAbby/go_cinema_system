package config

import (
	"fmt"
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

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

func Load() (Config, error) {
	var cfg Config
	var err error

	if err = godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные среды и значения по умолчанию.")
	}

	cfg.Server, err = LoadServerConfig()
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию сервера: %v", err)
		return Config{}, err
	}

	cfg.DB, err = LoadDBConfig()
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию базы данных: %v", err)
		return Config{}, err
	}

	return cfg, nil
}

// Загрузка конфига сервера
func LoadServerConfig() (ServerConfig, error) {
	var cfg ServerConfig
	var err error

	cfg.Addr, err = getEnv("APP_ADDR")
	if err != nil {
		return ServerConfig{}, err
	}
	cfg.ReadTimeout, err = getDurationEnv("APP_READ_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}
	cfg.WriteTimeout, err = getDurationEnv("APP_WRITE_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}
	cfg.ShutdownTimeout, err = getDurationEnv("APP_SHUTDOWN_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	return cfg, nil
}

// Загрузка конфига базы данных
func LoadDBConfig() (DBConfig, error) {
	var cfg DBConfig
	var err error

	cfg.Host, err = getEnv("DB_HOST")
	if err != nil {
		return DBConfig{}, err
	}
	cfg.Port, err = getEnv("DB_PORT")
	if err != nil {
		return DBConfig{}, err
	}
	cfg.User, err = getEnv("DB_USER")
	if err != nil {
		return DBConfig{}, err
	}
	cfg.Password, err = getEnv("DB_PASSWORD")
	if err != nil {
		return DBConfig{}, err
	}
	cfg.Name, err = getEnv("DB_NAME")
	if err != nil {
		return DBConfig{}, err
	}
	cfg.SSLMode, err = getEnv("DB_SSL_MODE")
	if err != nil {
		return DBConfig{}, err
	}

	return cfg, nil
}

// Функции получения из env
func getEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return v, nil
}

func getDurationEnv(key string) (time.Duration, error) {
	value, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при парсинге длительности.: %w", err)
	}

	return duration, nil
}
