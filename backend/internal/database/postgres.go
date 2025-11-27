package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"crypto-orderbook/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	connString := cfg.ConnectionString()

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Set pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
	log.Println("Database connection closed")
}

// RunMigrations runs SQL migration files
func (db *Database) RunMigrations() error {
	ctx := context.Background()

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,

		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			order_type VARCHAR(4) NOT NULL,
			price DECIMAL(18,8) NOT NULL CHECK (price > 0),
			amount DECIMAL(18,8) NOT NULL CHECK (amount > 0),
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT check_order_type CHECK (order_type IN ('buy', 'sell')),
			CONSTRAINT check_status CHECK (status IN ('active', 'filled', 'cancelled'))
		);
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
		CREATE INDEX IF NOT EXISTS idx_orders_type ON orders(order_type);
		CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);`,
	}

	for i, migration := range migrations {
		if _, err := db.Pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
		log.Printf("✅ Migration %d completed", i+1)
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}
