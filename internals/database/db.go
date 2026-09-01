package database

import (
	"context"
	"fmt"
	"os"
	"time"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB() (*DB, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || username == "" || password == "" || dbname == "" {
		log.Fatal("Missing one or more required environment variables")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbname)

	
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf("Error in connection %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{
		Pool: pool,
	}, nil
}
