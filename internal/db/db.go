package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create pool
func CreateDBPool(dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	return pool, nil
}

// Close pool
func CloseDBPool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
