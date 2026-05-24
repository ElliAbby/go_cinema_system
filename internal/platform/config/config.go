package config

import (
	"fmt"
	"log"
	"os"
	"strings"
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

type JWTConfig struct {
	SecretKey string
}

type KafkaConfig struct {
	Brokers              []string
	BookingPaymentsTopic  string
	BookingPaymentsGroup  string
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

type Config struct {
	Server                  ServerConfig
	DB                      DBConfig
	JWT                     JWTConfig
	Kafka                   KafkaConfig
	ReservationHoldDuration time.Duration
}

type WorkerConfig struct {
	DB          DBConfig
	Kafka       KafkaConfig
	ShutdownTimeout time.Duration
	MetricsAddr string
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
	cfg.JWT, err = LoadJWTConfig()
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию JWT: %v", err)
		return Config{}, err
	}
	cfg.Kafka, err = LoadKafkaConfig()
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию Kafka: %v", err)
		return Config{}, err
	}
	cfg.ReservationHoldDuration, err = getDurationEnvOrDefault("RESERVATION_HOLD_DURATION", 1*time.Minute)
	if err != nil {
		log.Printf("Не удалось загрузить длительность резерва: %v", err)
		return Config{}, err
	}
	return cfg, nil
}

func LoadWorkerConfig() (WorkerConfig, error) {
	var cfg WorkerConfig
	var err error

	cfg.DB, err = LoadDBConfig()
	if err != nil {
		return WorkerConfig{}, err
	}

	cfg.Kafka, err = LoadKafkaConfig()
	if err != nil {
		return WorkerConfig{}, err
	}
	cfg.ShutdownTimeout, err = getDurationEnv("APP_SHUTDOWN_TIMEOUT")
	if err != nil {
		return WorkerConfig{}, err
	}

	cfg.MetricsAddr = getEnvOrDefault("APP_METRICS_ADDR", ":8081")

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

// Загрузка конфига JWT
func LoadJWTConfig() (JWTConfig, error) {
	var cfg JWTConfig
	var err error
	
	cfg.SecretKey, err = getEnv("JWT_SECRET_KEY")
	if err != nil {
		return JWTConfig{}, err
	}
	return cfg, nil
}

// Загрузка конфига Kafka
func LoadKafkaConfig() (KafkaConfig, error) {
	var cfg KafkaConfig
	var err error

	cfg.Brokers, err = getCSVEnv("KAFKA_BROKERS")
	if err != nil {
		return KafkaConfig{}, err
	}
	cfg.BookingPaymentsTopic, err = getEnv("KAFKA_BOOKING_PAYMENTS_TOPIC")
	if err != nil {
		return KafkaConfig{}, err
	}
	cfg.BookingPaymentsGroup, err = getEnv("KAFKA_BOOKING_PAYMENTS_GROUP")
	if err != nil {
		return KafkaConfig{}, err
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

func getEnvOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
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

func getDurationEnvOrDefault(key string, defaultValue time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при парсинге длительности.: %w", err)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("environment variable %s must be greater than zero", key)
	}

	return duration, nil
}

func getCSVEnv(key string) ([]string, error) {
	value, err := getEnv(key)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("environment variable %s is empty", key)
	}

	return items, nil
}
