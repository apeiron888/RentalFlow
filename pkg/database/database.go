package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rentalflow/rentalflow/pkg/config"
	"github.com/rentalflow/rentalflow/pkg/logger"
)

// DB holds the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection
func New(cfg config.DatabaseConfig) (*DB, error) {
	log := logger.NewLogger("database")

	// Build connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_min_conns=%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
		cfg.MaxOpenConns,
		cfg.MaxIdleConns,
	)

	// Parse config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConnLifetime = cfg.MaxLifetime

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Msg("Connected to database")

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Health checks if the database is healthy
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
