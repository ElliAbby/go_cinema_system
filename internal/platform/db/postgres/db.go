package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ElliAbby/go_cinema_system/internal/platform/config"

)

const (
	driverName = "postgres"
)

func NewPostgresDB(cfg *config.DBConfig) (*sqlx.DB, error) {
	dsn := cfg.DSN()

	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func IsDatabaseEmpty(db *sqlx.DB) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM movies")
	if err != nil {
		return true, err
	}

	return count == 0, nil
}

func SeedDatabase(db *sqlx.DB) error {
	isEmpty, err := IsDatabaseEmpty(db)
	if err != nil {
		return fmt.Errorf("failed to check if database is empty: %w", err)
	}

	if !isEmpty {
		return nil
	}

	seedPath := "deployment/db/init/02_seed_data.sql"
	
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
		seedPath = filepath.Join(cwd, "deployment/db/init/02_seed_data.sql")
		
		if _, err := os.Stat(seedPath); os.IsNotExist(err) {
			return fmt.Errorf("seed file not found at %s", seedPath)
		}
	}

	seedSQL, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, string(seedSQL))
	if err != nil {
		return fmt.Errorf("failed to execute seed script: %w", err)
	}

	return nil
}